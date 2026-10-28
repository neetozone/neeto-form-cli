package commands

import (
	"net/http"

	"github.com/neetozone/neeto-cli-commons/output"
	"github.com/spf13/cobra"
)

const questionAttributesHelp = `Path to a JSON file with the question attributes, e.g. {"label": "Work email", "is_required": true, "options": [{"label": "Basic"}]}. Run "neetoform question-kinds list" for the attributes of each kind`

var questionsCmd = &cobra.Command{
	Use:   "questions",
	Short: "Manage form questions",
}

func questionAttributes(cmd *cobra.Command) (map[string]any, error) {
	return jsonFileWithFlags(cmd, "label", "description", "placeholder", "required=is_required")
}

func addQuestionAttributeFlags(cmd *cobra.Command) {
	addJSONFileFlag(cmd, questionAttributesHelp)
	cmd.Flags().String("label", "", "Question label")
	cmd.Flags().String("description", "", "Question description")
	cmd.Flags().String("placeholder", "", "Input placeholder")
	cmd.Flags().Bool("required", false, "Whether an answer is required")
}

var questionsCreateCmd = &cobra.Command{
	Use:   "create <form-id>",
	Short: "Create a question",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		attributes, err := questionAttributes(cmd)
		if err != nil {
			return err
		}

		question := flagValues(cmd, "kind", "position")
		if len(attributes) > 0 {
			question["attributes"] = attributes
		}

		return sendAndPrint(cmd, http.MethodPost, formPath(args[0], "questions"), map[string]any{"question": question}, printActionResult)
	},
}

var questionsBulkCreateCmd = &cobra.Command{
	Use:   "bulk-create <form-id>",
	Short: "Create several questions at once",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		body, err := jsonFileBody(cmd)
		if err != nil {
			return err
		}
		return sendAndPrint(cmd, http.MethodPost, formPath(args[0], "questions", "bulk_create"), body, printActionResult)
	},
}

var questionsUpdateCmd = &cobra.Command{
	Use:   "update <form-id> <question-id>",
	Short: "Update a question",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		attributes, err := questionAttributes(cmd)
		if err != nil {
			return err
		}

		body := map[string]any{"question": map[string]any{"attributes": attributes}}
		return sendAndPrint(cmd, http.MethodPatch, formPath(args[0], "questions", args[1]), body, printResource)
	},
}

var questionsCloneCmd = &cobra.Command{
	Use:   "clone <form-id> <question-id>",
	Short: "Clone a question",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		return sendAndPrint(cmd, http.MethodPost, formPath(args[0], "questions", args[1], "clone"), nil, printActionResult)
	},
}

var questionsReorderCmd = &cobra.Command{
	Use:   "reorder <form-id>",
	Short: "Reorder the questions of a form",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		body := flagValues(cmd, "ids[]=question_ids")
		return sendAndPrint(cmd, http.MethodPatch, formPath(args[0], "questions", "reorder"), body, printActionResult)
	},
}

func init() {
	questionsCreateCmd.Flags().String("kind", "", "Question kind, e.g. short_text or dropdown")
	questionsCreateCmd.Flags().Int("position", 0, "Zero-based position; defaults to just before the submit button")
	addQuestionAttributeFlags(questionsCreateCmd)
	markFlagsRequired(questionsCreateCmd, "kind")

	addJSONFileFlag(questionsBulkCreateCmd, `Path to a JSON file, e.g. {"questions": [{"kind": "phone"}, {"kind": "yes_no", "attributes": {"label": "Subscribe?"}}]}`)
	markFlagsRequired(questionsBulkCreateCmd, "json-file")

	addQuestionAttributeFlags(questionsUpdateCmd)

	questionsReorderCmd.Flags().String("ids", "", "Comma-separated ids of every element of the form in the new order, logo, welcome screen and submit button included")
	markFlagsRequired(questionsReorderCmd, "ids")

	questionsCmd.AddCommand(
		listCmd("list <form-id>", "List questions", 1, formResource("questions"), "questions", []output.Breadcrumb{
			{Label: "Show", Command: "neetoform forms questions show <form-id> <question-id>"},
			{Label: "Kinds", Command: "neetoform question-kinds list"},
		}),
		showCmd("show <form-id> <question-id>", "Show a question", 2, formResource("questions")),
		questionsCreateCmd,
		questionsBulkCreateCmd,
		questionsUpdateCmd,
		deleteCmd("delete <form-id> <question-id>", "Delete a question", 2, formResource("questions"), "Question deleted."),
		questionsCloneCmd,
		questionsReorderCmd,
	)
	formsCmd.AddCommand(questionsCmd)
}
