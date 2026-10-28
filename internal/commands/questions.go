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
	attributes, err := jsonFileBody(cmd)
	if err != nil {
		return nil, err
	}
	for key, value := range flagValues(cmd, "label", "description", "placeholder", "required=is_required") {
		attributes[key] = value
	}
	return attributes, nil
}

func addQuestionAttributeFlags(cmd *cobra.Command) {
	cmd.Flags().String("json-file", "", questionAttributesHelp)
	cmd.Flags().String("label", "", "Question label")
	cmd.Flags().String("description", "", "Question description")
	cmd.Flags().String("placeholder", "", "Input placeholder")
	cmd.Flags().Bool("required", false, "Whether an answer is required")
}

var questionsListCmd = &cobra.Command{
	Use:   "list <form-id>",
	Short: "List questions",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := fetch(cmd, formPath(args[0], "questions"), paginationParams(cmd))
		if err != nil {
			return err
		}

		printList(data, "questions", []output.Breadcrumb{
			{Label: "Show", Command: "neetoform forms questions show <form-id> <question-id>"},
			{Label: "Kinds", Command: "neetoform question-kinds list"},
		})
		return nil
	},
}

var questionsShowCmd = &cobra.Command{
	Use:   "show <form-id> <question-id>",
	Short: "Show a question",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := fetch(cmd, formPath(args[0], "questions", args[1]), nil)
		if err != nil {
			return err
		}

		printResource(data, nil)
		return nil
	},
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

		data, err := send(cmd, http.MethodPost, formPath(args[0], "questions"), map[string]any{"question": question})
		if err != nil {
			return err
		}

		printActionResult(data, nil)
		return nil
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

		data, err := send(cmd, http.MethodPost, formPath(args[0], "questions", "bulk_create"), body)
		if err != nil {
			return err
		}

		printActionResult(data, nil)
		return nil
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
		data, err := send(cmd, http.MethodPatch, formPath(args[0], "questions", args[1]), body)
		if err != nil {
			return err
		}

		printResource(data, nil)
		return nil
	},
}

var questionsDeleteCmd = &cobra.Command{
	Use:   "delete <form-id> <question-id>",
	Short: "Delete a question",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		return removeAndConfirm(cmd, formPath(args[0], "questions", args[1]), nil, "Question deleted.")
	},
}

var questionsCloneCmd = &cobra.Command{
	Use:   "clone <form-id> <question-id>",
	Short: "Clone a question",
	Args:  cobra.ExactArgs(2),
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := send(cmd, http.MethodPost, formPath(args[0], "questions", args[1], "clone"), nil)
		if err != nil {
			return err
		}

		printActionResult(data, nil)
		return nil
	},
}

var questionsReorderCmd = &cobra.Command{
	Use:   "reorder <form-id>",
	Short: "Reorder the questions of a form",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		body := flagValues(cmd, "ids[]=question_ids")
		data, err := send(cmd, http.MethodPatch, formPath(args[0], "questions", "reorder"), body)
		if err != nil {
			return err
		}

		printActionResult(data, nil)
		return nil
	},
}

func init() {
	addPaginationFlags(questionsListCmd)

	questionsCreateCmd.Flags().String("kind", "", "Question kind, e.g. short_text or dropdown")
	questionsCreateCmd.Flags().Int("position", 0, "Zero-based position; defaults to just before the submit button")
	addQuestionAttributeFlags(questionsCreateCmd)
	markFlagsRequired(questionsCreateCmd, "kind")

	questionsBulkCreateCmd.Flags().String("json-file", "", `Path to a JSON file, e.g. {"questions": [{"kind": "phone"}, {"kind": "yes_no", "attributes": {"label": "Subscribe?"}}]}`)
	markFlagsRequired(questionsBulkCreateCmd, "json-file")

	addQuestionAttributeFlags(questionsUpdateCmd)

	questionsReorderCmd.Flags().String("ids", "", "Comma-separated ids of every element of the form in the new order, logo, welcome screen and submit button included")
	markFlagsRequired(questionsReorderCmd, "ids")

	questionsCmd.AddCommand(questionsListCmd, questionsShowCmd, questionsCreateCmd, questionsBulkCreateCmd,
		questionsUpdateCmd, questionsDeleteCmd, questionsCloneCmd, questionsReorderCmd)
	formsCmd.AddCommand(questionsCmd)
}
