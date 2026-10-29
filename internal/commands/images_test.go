package commands

import (
	"crypto/md5"
	"encoding/base64"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeImage(t *testing.T, name string, data []byte) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestImagesUploadReservesThenPutsBytesAndHidesCredentials(t *testing.T) {
	var putMethod, putPath string
	var putHeaders http.Header
	var putBody []byte
	storage := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		putMethod, putPath, putHeaders = r.Method, r.URL.Path, r.Header.Clone()
		putBody, _ = io.ReadAll(r.Body)
	}))
	t.Cleanup(storage.Close)

	data := []byte("\x89PNG-bytes")
	sum := md5.Sum(data)
	checksum := base64.StdEncoding.EncodeToString(sum[:])
	response := fmt.Sprintf(`{"signed_id":"sig123","filename":"logo.png","content_type":"image/png","byte_size":%d,"checksum":%q,"upload_url":"%s/bucket/secret-key","upload_headers":{"Content-Type":"image/png","Content-MD5":%q}}`,
		len(data), checksum, storage.URL, checksum)

	got, out, err := runCLI(t, response, "images", "upload", writeImage(t, "logo.png", data))
	if err != nil {
		t.Fatal(err)
	}
	if got.method != "POST" || got.path != "/api/external/v2/image_uploads" {
		t.Errorf("request = %s %s", got.method, got.path)
	}
	if got.body["filename"] != "logo.png" || got.body["byte_size"] != float64(len(data)) ||
		got.body["checksum"] != checksum || got.body["content_type"] != "image/png" {
		t.Errorf("body = %v", got.body)
	}
	if putMethod != "PUT" || putPath != "/bucket/secret-key" || string(putBody) != string(data) {
		t.Errorf("put = %s %s %q", putMethod, putPath, putBody)
	}
	if putHeaders.Get("Content-Type") != "image/png" || putHeaders.Get("Content-MD5") != checksum || putHeaders.Get("Authorization") != "" {
		t.Errorf("put headers = %v", putHeaders)
	}
	if !strings.Contains(out, "sig123") || strings.Contains(out, "secret-key") || strings.Contains(out, "upload_") || strings.Contains(out, "Content-MD5") {
		t.Errorf("output = %q", out)
	}
}

func TestImagesUploadRefusesFileOverOneMegabyteWithoutRequest(t *testing.T) {
	path := writeImage(t, "big.png", make([]byte, 1<<20+1))

	got, _, err := runCLI(t, `{}`, "images", "upload", path)
	if err == nil || !strings.Contains(err.Error(), "larger than 1 MB") {
		t.Fatalf("error = %v", err)
	}
	if got.path != "" {
		t.Errorf("request was sent to %q", got.path)
	}
}

func TestImagesUploadRefusesUnsupportedTypeWithoutRequest(t *testing.T) {
	got, _, err := runCLI(t, `{}`, "images", "upload", writeImage(t, "notes.txt", []byte("hi")))
	if err == nil || !strings.Contains(err.Error(), "Unsupported image type") {
		t.Fatalf("error = %v", err)
	}
	if got.path != "" {
		t.Errorf("request was sent to %q", got.path)
	}
}

func TestImagesUploadRefusesMissingFileWithoutRequest(t *testing.T) {
	got, _, err := runCLI(t, `{}`, "images", "upload", filepath.Join(t.TempDir(), "missing.png"))
	if err == nil || !strings.Contains(err.Error(), "file not found") {
		t.Fatalf("error = %v", err)
	}
	if got.path != "" {
		t.Errorf("request was sent to %q", got.path)
	}
}

func TestImagesCreateSendsUrlAndFilename(t *testing.T) {
	got, out, err := runCLI(t, `{"signed_id":"sig9","filename":"hero.png","content_type":"image/png","byte_size":10}`,
		"images", "create", "--url", "https://example.com/hero.png", "--filename", "hero.png")
	if err != nil {
		t.Fatal(err)
	}
	if got.method != "POST" || got.path != "/api/external/v2/images" {
		t.Errorf("request = %s %s", got.method, got.path)
	}
	if got.body["url"] != "https://example.com/hero.png" || got.body["filename"] != "hero.png" {
		t.Errorf("body = %v", got.body)
	}
	if !strings.Contains(out, "sig9") {
		t.Errorf("output = %q", out)
	}
}
