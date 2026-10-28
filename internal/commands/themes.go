package commands

import (
	"net/http"

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
	theme, err := jsonFileWithFlags(cmd, "name", "base-theme-id")
	if err != nil {
		return nil, err
	}
	pairs, _ := cmd.Flags().GetStringArray("property")
	parsed, err := parseKeyValues(pairs, "property", "key=value")
	if err != nil {
		return nil, err
	}
	properties := map[string]any{}
	for _, pair := range parsed {
		properties[pair.key] = pair.value
	}
	mergeInto(theme, "properties", properties)
	return map[string]any{"theme": theme}, nil
}

func saveTheme(method string, render printer, path argsPath) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		body, err := themeBody(cmd)
		if err != nil {
			return err
		}
		return sendAndPrint(cmd, method, path(args), body, render)
	}
}

var themesCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a theme",
	RunE:  saveTheme(http.MethodPost, printActionResult, func([]string) string { return apiPath("themes") }),
}

var themesUpdateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update a theme",
	Args:  cobra.ExactArgs(1),
	RunE:  saveTheme(http.MethodPatch, printResource, workspaceResource("themes")),
}

var formThemeUpdateCmd = &cobra.Command{
	Use:   "update <form-id>",
	Short: "Apply a theme to a form",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return sendAndPrint(cmd, http.MethodPut, formPath(args[0], "theme"), flagValues(cmd, "theme-id"), printResource)
	},
}

func init() {
	propertyHelp := "Theme property as key=value, repeatable; e.g. --property background_color=#ffffff. Editable keys: eui_font_family, question_color, background_color, button_text_color, button_background_color, button_border_radius, welcome_text_alignment, questions_text_alignment, thank_you_text_alignment (left, center or right)"
	jsonFileHelp := "Path to a JSON file with the theme fields; e.g. {\"name\": \"Brand\", \"properties\": {\"question_color\": \"#123456\"}}. Flags override the file"

	themesCreateCmd.Flags().String("name", "", "Theme name")
	themesCreateCmd.Flags().String("base-theme-id", "", "Theme to copy; the default theme when omitted")
	themesCreateCmd.Flags().StringArray("property", nil, propertyHelp)
	addJSONFileFlag(themesCreateCmd, jsonFileHelp)

	themesUpdateCmd.Flags().String("name", "", "New theme name")
	themesUpdateCmd.Flags().StringArray("property", nil, propertyHelp)
	addJSONFileFlag(themesUpdateCmd, jsonFileHelp)

	formThemeUpdateCmd.Flags().String("theme-id", "", "ID of the system or workspace theme to apply")
	markFlagsRequired(formThemeUpdateCmd, "theme-id")

	themesCmd.AddCommand(
		listCmd("list", "List themes", 0, workspaceResource("themes"), "themes", []output.Breadcrumb{
			{Label: "Show", Command: "neetoform themes show <id>"},
		}),
		showCmd("show <id>", "Show a theme", 1, workspaceResource("themes")),
		themesCreateCmd,
		themesUpdateCmd,
		deleteCmd("delete <id>", "Delete a theme", 1, workspaceResource("themes"), "Theme deleted."),
	)
	formThemeCmd.AddCommand(formThemeUpdateCmd)
	formsCmd.AddCommand(formThemeCmd)
	register(func(root *cobra.Command) { root.AddCommand(themesCmd) })
}
