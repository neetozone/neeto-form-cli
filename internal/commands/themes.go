package commands

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/neetozone/neeto-cli-commons/output"
	"github.com/spf13/cobra"
)

var themesCmd = &cobra.Command{
	Use:   "themes",
	Short: "Manage themes",
}

var formThemeCmd = &cobra.Command{
	Use:   "theme",
	Short: "Manage a form's theme",
}

func themeBody(cmd *cobra.Command) (map[string]any, error) {
	properties, _ := cmd.Flags().GetStringArray("property")
	theme, err := jsonFileBody(cmd)
	if err != nil {
		return nil, err
	}
	for key, value := range flagValues(cmd, "name", "base-theme-id") {
		theme[key] = value
	}

	if len(properties) > 0 {
		merged, _ := theme["properties"].(map[string]any)
		if merged == nil {
			merged = map[string]any{}
		}
		for _, property := range properties {
			key, value, found := strings.Cut(property, "=")
			if !found || key == "" {
				return nil, fmt.Errorf("Invalid --property %q. Use key=value.", property)
			}
			merged[key] = value
		}
		theme["properties"] = merged
	}
	return map[string]any{"theme": theme}, nil
}

var themesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List themes",
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := fetch(cmd, apiPath("themes"), paginationParams(cmd))
		if err != nil {
			return err
		}
		printList(data, "themes", []output.Breadcrumb{
			{Label: "Show", Command: "neetoform themes show <id>"},
		})
		return nil
	},
}

var themesShowCmd = &cobra.Command{
	Use:   "show <id>",
	Short: "Show a theme",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := fetch(cmd, apiPath("themes", args[0]), nil)
		if err != nil {
			return err
		}
		printResource(data, nil)
		return nil
	},
}

var themesCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a theme",
	RunE: func(cmd *cobra.Command, args []string) error {
		body, err := themeBody(cmd)
		if err != nil {
			return err
		}
		data, err := send(cmd, http.MethodPost, apiPath("themes"), body)
		if err != nil {
			return err
		}
		printActionResult(data, nil)
		return nil
	},
}

var themesUpdateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update a theme",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		body, err := themeBody(cmd)
		if err != nil {
			return err
		}
		data, err := send(cmd, http.MethodPatch, apiPath("themes", args[0]), body)
		if err != nil {
			return err
		}
		printResource(data, nil)
		return nil
	},
}

var themesDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a theme",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return removeAndConfirm(cmd, apiPath("themes", args[0]), nil, "Theme deleted.")
	},
}

var formThemeUpdateCmd = &cobra.Command{
	Use:   "update <form-id>",
	Short: "Apply a theme to a form",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := send(cmd, http.MethodPut, formPath(args[0], "theme"), flagValues(cmd, "theme-id"))
		if err != nil {
			return err
		}
		printResource(data, nil)
		return nil
	},
}

func init() {
	addPaginationFlags(themesListCmd)

	propertyHelp := "Theme property as key=value, repeatable; e.g. --property background_color=#ffffff. Editable keys: eui_font_family, question_color, background_color, button_text_color, button_background_color, button_border_radius, welcome_text_alignment, questions_text_alignment, thank_you_text_alignment (left, center or right)"
	jsonFileHelp := "Path to a JSON file with the theme fields; e.g. {\"name\": \"Brand\", \"properties\": {\"question_color\": \"#123456\"}}. Flags override the file"

	themesCreateCmd.Flags().String("name", "", "Theme name")
	themesCreateCmd.Flags().String("base-theme-id", "", "Theme to copy; the default theme when omitted")
	themesCreateCmd.Flags().StringArray("property", nil, propertyHelp)
	themesCreateCmd.Flags().String("json-file", "", jsonFileHelp)

	themesUpdateCmd.Flags().String("name", "", "New theme name")
	themesUpdateCmd.Flags().StringArray("property", nil, propertyHelp)
	themesUpdateCmd.Flags().String("json-file", "", jsonFileHelp)

	formThemeUpdateCmd.Flags().String("theme-id", "", "ID of the system or workspace theme to apply")
	markFlagsRequired(formThemeUpdateCmd, "theme-id")

	themesCmd.AddCommand(themesListCmd, themesShowCmd, themesCreateCmd, themesUpdateCmd, themesDeleteCmd)
	formThemeCmd.AddCommand(formThemeUpdateCmd)
	formsCmd.AddCommand(formThemeCmd)
	register(func(root *cobra.Command) { root.AddCommand(themesCmd) })
}
