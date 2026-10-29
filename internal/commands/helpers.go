package commands

import (
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"net/http"
	"net/url"
	"slices"
	"strings"

	"github.com/neetozone/neeto-cli-commons/output"
	"github.com/spf13/cobra"
)

type argsPath func(args []string) string

type printer func(data json.RawMessage, breadcrumbs []output.Breadcrumb)

type keyValue struct{ key, value string }

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

func workspaceResource(segment string) argsPath {
	return func(args []string) string {
		return apiPath(append([]string{segment}, args...)...)
	}
}

func formResource(segment string) argsPath {
	return func(args []string) string {
		return formPath(args[0], append([]string{segment}, args[1:]...)...)
	}
}

func showCmd(use, short string, nargs int, path argsPath, querySpecs ...string) *cobra.Command {
	return &cobra.Command{
		Use:   use,
		Short: short,
		Args:  cobra.ExactArgs(nargs),
		RunE: func(cmd *cobra.Command, args []string) error {
			data, err := fetch(cmd, path(args), queryValues(cmd, url.Values{}, querySpecs...))
			if err != nil {
				return err
			}
			printResource(data, nil)
			return nil
		},
	}
}

func listCmd(use, short string, nargs int, path argsPath, key string, breadcrumbs []output.Breadcrumb, querySpecs ...string) *cobra.Command {
	list := &cobra.Command{
		Use:   use,
		Short: short,
		Args:  cobra.ExactArgs(nargs),
		RunE: func(cmd *cobra.Command, args []string) error {
			data, err := fetch(cmd, path(args), queryValues(cmd, paginationParams(cmd), querySpecs...))
			if err != nil {
				return err
			}
			printList(data, key, breadcrumbs)
			return nil
		},
	}
	addPaginationFlags(list)
	return list
}

func saveCmd(use, short string, nargs int, method string, path argsPath, render printer, build func(cmd *cobra.Command) (map[string]any, error)) *cobra.Command {
	return &cobra.Command{
		Use:   use,
		Short: short,
		Args:  cobra.ExactArgs(nargs),
		RunE: func(cmd *cobra.Command, args []string) error {
			body, err := build(cmd)
			if err != nil {
				return err
			}
			if hasEmptyResource(body) {
				return errors.New("Nothing to send. Pass --json-file or set at least one flag; see --help.")
			}
			return sendAndPrint(cmd, method, path(args), body, render)
		},
	}
}

func hasEmptyResource(body map[string]any) bool {
	for _, value := range body {
		if resource, ok := value.(map[string]any); ok && len(resource) == 0 {
			return true
		}
	}
	return false
}

func deleteCmd(use, short string, nargs int, path argsPath, message string) *cobra.Command {
	return &cobra.Command{
		Use:   use,
		Short: short,
		Args:  cobra.ExactArgs(nargs),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := getClient(cmd)
			if err != nil {
				return err
			}
			if err := c.Delete(path(args)); err != nil {
				return err
			}
			printMessage(message)
			return nil
		},
	}
}

func addJSONFileFlag(cmd *cobra.Command, help string) {
	cmd.Flags().String("json-file", "", help)
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

func validateEnumList(flag string, values, valid []string) error {
	for _, value := range values {
		if err := validateEnum(flag, value, valid); err != nil {
			return err
		}
	}
	return nil
}

func parseKeyValues(pairs []string, flag, usage string) ([]keyValue, error) {
	parsed := make([]keyValue, 0, len(pairs))
	for _, pair := range pairs {
		key, value, found := strings.Cut(pair, "=")
		if !found || key == "" {
			return nil, fmt.Errorf("Invalid --%s %q. Use %s.", flag, pair, usage)
		}
		parsed = append(parsed, keyValue{key, value})
	}
	return parsed, nil
}

func mergeInto(body map[string]any, key string, extra map[string]any) {
	if len(extra) == 0 {
		return
	}
	merged, _ := body[key].(map[string]any)
	if merged == nil {
		merged = map[string]any{}
	}
	maps.Copy(merged, extra)
	body[key] = merged
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
	path := flagString(cmd, "json-file")
	if path == "" {
		return map[string]any{}, nil
	}
	return readJSONFile(path)
}

func jsonFileWithFlags(cmd *cobra.Command, specs ...string) (map[string]any, error) {
	body, err := jsonFileBody(cmd)
	if err != nil {
		return nil, err
	}
	maps.Copy(body, flagValues(cmd, specs...))
	return body, nil
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

func sendAndPrint(cmd *cobra.Command, method, path string, body any, render printer) error {
	data, err := send(cmd, method, path, body)
	if err != nil {
		return err
	}
	render(data, nil)
	return nil
}
