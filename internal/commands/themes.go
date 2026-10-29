package commands

import (
	"errors"
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
	theme, err := jsonFileWithFlags(cmd, "name", "base-theme-id", "background-image")
	if err != nil {
		return nil, err
	}
	if cmd.Flags().Changed("background-image") && cmd.Flags().Changed("remove-background-image") {
		return nil, errors.New("Use either --background-image or --remove-background-image, not both.")
	}
	if cmd.Flags().Changed("remove-background-image") {
		theme["background_image"] = nil
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

var themesCreateCmd = saveCmd("create", "Create a theme", 0, http.MethodPost, workspaceResource("themes"), printActionResult, themeBody)

var themesUpdateCmd = saveCmd("update <id>", "Update a theme", 1, http.MethodPatch, workspaceResource("themes"), printResource, themeBody)

var formThemeUpdateCmd = &cobra.Command{
	Use:   "update <form-id>",
	Short: "Apply a theme to a form",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return sendAndPrint(cmd, http.MethodPut, formPath(args[0], "theme"), flagValues(cmd, "theme-id"), printResource)
	},
}

func addImageFlags(cmd *cobra.Command) {
	cmd.Flags().String("background-image", "", "signed_id from `neetoform images upload` or `images create`")
	cmd.Flags().Bool("remove-background-image", false, "Remove the background image")
}

func init() {
	propertyHelp := "Theme property as key=value, repeatable; e.g. --property background_color=#ffffff. Editable keys: eui_font_family, question_color, answer_color, answer_text_color, background_color, button_text_color, button_background_color, button_border_radius, welcome_text_alignment, questions_text_alignment, thank_you_text_alignment (left, center or right), is_content_background_enabled (true or false), content_background_color (hex), background_image_position (center, top left, top center, top right, center left, center right, bottom left, bottom right or bottom center), background (image overlay, an integer from -100 to 100), is_dominant_background_color_enabled (true or false), dominant_color (hex), custom_css (CSS added inside the form; start each selector with .neeto-form-eui and target elements by the [id] selectors in `forms outline`; Pro plan only; an empty value removes it). Values are sent as given and checked by the server"
	jsonFileHelp := "Path to a JSON file with the theme fields; e.g. {\"name\": \"Brand\", \"properties\": {\"question_color\": \"#123456\"}}. Flags override the file"

	themesCreateCmd.Flags().String("name", "", "Theme name")
	themesCreateCmd.Flags().String("base-theme-id", "", "Theme to copy; the default theme when omitted")
	themesCreateCmd.Flags().StringArray("property", nil, propertyHelp)
	addImageFlags(themesCreateCmd)
	addJSONFileFlag(themesCreateCmd, jsonFileHelp)

	themesUpdateCmd.Flags().String("name", "", "New theme name")
	themesUpdateCmd.Flags().StringArray("property", nil, propertyHelp)
	addImageFlags(themesUpdateCmd)
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
