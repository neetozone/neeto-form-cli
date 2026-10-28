package commands

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"strings"

	"github.com/spf13/cobra"
)

func apiPath(segments ...string) string {
	escaped := make([]string, len(segments))
	for i, segment := range segments {
		escaped[i] = url.PathEscape(segment)
	}
	return "/" + strings.Join(escaped, "/")
}

func formPath(formID string, segments ...string) string {
	return apiPath(append([]string{"forms", formID}, segments...)...)
}

func flagString(cmd *cobra.Command, name string) string {
	value, _ := cmd.Flags().GetString(name)
	return value
}

func validateEnum(flag, value string, valid []string) error {
	if value == "" || slices.Contains(valid, value) {
		return nil
	}
	return fmt.Errorf("Invalid --%s %q. Valid values: %s.", flag, value, strings.Join(valid, ", "))
}

func flagValues(cmd *cobra.Command, specs ...string) map[string]interface{} {
	values := map[string]interface{}{}
	for _, spec := range specs {
		name, key, _ := strings.Cut(spec, "=")
		csv := strings.HasSuffix(name, "[]")
		name = strings.TrimSuffix(name, "[]")
		if key == "" {
			key = strings.ReplaceAll(name, "-", "_")
		}

		flag := cmd.Flags().Lookup(name)
		if flag == nil || !flag.Changed {
			continue
		}

		switch {
		case csv:
			values[key] = splitCSV(flag.Value.String())
		case flag.Value.Type() == "bool":
			values[key], _ = cmd.Flags().GetBool(name)
		case flag.Value.Type() == "int":
			values[key], _ = cmd.Flags().GetInt(name)
		case flag.Value.Type() == "stringArray":
			values[key], _ = cmd.Flags().GetStringArray(name)
		default:
			values[key] = flag.Value.String()
		}
	}
	return values
}

func queryValues(cmd *cobra.Command, params url.Values, specs ...string) url.Values {
	for key, value := range flagValues(cmd, specs...) {
		if list, ok := value.([]string); ok {
			for _, item := range list {
				params.Add(key+"[]", item)
			}
			continue
		}
		params.Set(key, fmt.Sprint(value))
	}
	return params
}

func jsonFileBody(cmd *cobra.Command) (map[string]any, error) {
	path, _ := cmd.Flags().GetString("json-file")
	if path == "" {
		return map[string]any{}, nil
	}
	return readJSONFile(path)
}

func fetch(cmd *cobra.Command, path string, params url.Values) (json.RawMessage, error) {
	c, err := getClient(cmd)
	if err != nil {
		return nil, err
	}
	return c.Get(path, params)
}

func send(cmd *cobra.Command, method, path string, body any) (json.RawMessage, error) {
	c, err := getClient(cmd)
	if err != nil {
		return nil, err
	}

	switch method {
	case http.MethodPost:
		return c.Post(path, body)
	case http.MethodPut:
		return c.Put(path, body)
	case http.MethodPatch:
		return c.Patch(path, body)
	case http.MethodDelete:
		return c.DeleteWithBody(path, body)
	default:
		return nil, fmt.Errorf("unsupported method %s", method)
	}
}

func removeAndConfirm(cmd *cobra.Command, path string, params url.Values, message string) error {
	c, err := getClient(cmd)
	if err != nil {
		return err
	}
	if err := c.DeleteWithParams(path, params); err != nil {
		return err
	}
	printMessage(message)
	return nil
}
