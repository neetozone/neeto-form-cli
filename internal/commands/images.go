package commands

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	neturl "net/url"
	"os"
	"path/filepath"
	"slices"
	"strings"

	"github.com/spf13/cobra"
)

const maxImageBytes = 1 << 20

var imageContentTypes = map[string]string{
	".png":  "image/png",
	".jpg":  "image/jpeg",
	".jpeg": "image/jpeg",
	".gif":  "image/gif",
	".svg":  "image/svg+xml",
}

var imagesCmd = &cobra.Command{
	Use:   "images",
	Short: "Upload images for themes",
}

type imageUpload struct {
	SignedID      string            `json:"signed_id"`
	Filename      string            `json:"filename"`
	ContentType   string            `json:"content_type"`
	ByteSize      int64             `json:"byte_size"`
	UploadURL     string            `json:"upload_url"`
	UploadHeaders map[string]string `json:"upload_headers"`
}

func imageContentType(path, override string) (string, error) {
	contentType := override
	if contentType == "" {
		contentType = imageContentTypes[strings.ToLower(filepath.Ext(path))]
	}
	allowed := []string{"image/png", "image/jpeg", "image/gif", "image/svg+xml"}
	if !slices.Contains(allowed, contentType) {
		return "", fmt.Errorf("Unsupported image type for %q. Use png, jpeg, gif or svg.", filepath.Base(path))
	}
	return contentType, nil
}

func putImageBytes(upload imageUpload, data []byte) error {
	req, err := http.NewRequest(http.MethodPut, upload.UploadURL, bytes.NewReader(data))
	if err != nil {
		return errors.New("The server returned an invalid upload address.")
	}
	for name, value := range upload.UploadHeaders {
		req.Header.Set(name, value)
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		var urlErr *neturl.Error
		if errors.As(err, &urlErr) {
			err = urlErr.Err
		}
		return fmt.Errorf("Could not upload the image bytes: %w", err)
	}
	defer func() { _ = resp.Body.Close() }()
	_, _ = io.Copy(io.Discard, resp.Body)
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf("Could not upload the image bytes: storage answered %s.", resp.Status)
	}
	return nil
}

var imagesUploadCmd = &cobra.Command{
	Use:   "upload <path>",
	Short: "Upload a local image file and print its signed_id",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		path := args[0]
		info, err := os.Stat(path)
		if err != nil {
			return fmt.Errorf("Cannot read %q: file not found.", path)
		}
		if info.IsDir() {
			return fmt.Errorf("Cannot read %q: it is a directory.", path)
		}
		if info.Size() > maxImageBytes {
			return fmt.Errorf("%q is larger than 1 MB. Use a smaller image.", path)
		}
		contentType, err := imageContentType(path, flagString(cmd, "content-type"))
		if err != nil {
			return err
		}
		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		sum := md5.Sum(data)
		body := map[string]any{
			"filename":     filepath.Base(path),
			"byte_size":    len(data),
			"checksum":     base64.StdEncoding.EncodeToString(sum[:]),
			"content_type": contentType,
		}
		response, err := send(cmd, http.MethodPost, apiPath("image_uploads"), body)
		if err != nil {
			return err
		}
		var upload imageUpload
		if err := json.Unmarshal(response, &upload); err != nil {
			return err
		}
		if err := putImageBytes(upload, data); err != nil {
			return err
		}
		result, err := json.Marshal(map[string]any{
			"signed_id":    upload.SignedID,
			"filename":     upload.Filename,
			"content_type": upload.ContentType,
			"byte_size":    upload.ByteSize,
		})
		if err != nil {
			return err
		}
		printActionResult(result, nil)
		return nil
	},
}

var imagesCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Store an image from a public https URL and print its signed_id",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		return sendAndPrint(cmd, http.MethodPost, apiPath("images"), flagValues(cmd, "url", "filename"), printActionResult)
	},
}

func init() {
	imagesUploadCmd.Flags().String("content-type", "", "Image type when the file extension does not show it: image/png, image/jpeg, image/gif or image/svg+xml")

	imagesCreateCmd.Flags().String("url", "", "Public https URL of the image (png, jpeg, gif or svg, at most 1 MB)")
	imagesCreateCmd.Flags().String("filename", "", "Name to store the image under; taken from the URL when omitted")
	markFlagsRequired(imagesCreateCmd, "url")

	imagesCmd.AddCommand(imagesUploadCmd, imagesCreateCmd)
	register(func(root *cobra.Command) { root.AddCommand(imagesCmd) })
}
