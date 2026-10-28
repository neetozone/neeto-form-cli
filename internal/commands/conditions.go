package commands

import (
	"net/http"
	"net/url"

	"github.com/neetozone/neeto-cli-commons/output"
	"github.com/spf13/cobra"
)

type logicResource struct {
	use      string
	path     string
	singular string
	listKey  string
	idsKey   string
	label    string
	search   bool
}

func logicPayload(cmd *cobra.Command, key string, specs ...string) (map[string]any, error) {
	body, err := jsonFileBody(cmd)
	if err != nil {
		return nil, err
	}
	for k, v := range flagValues(cmd, specs...) {
		body[k] = v
	}
	return map[string]any{key: body}, nil
}

func newLogicGroup(r logicResource) *cobra.Command {
	group := &cobra.Command{Use: r.use, Short: "Manage " + r.label}

	list := &cobra.Command{
		Use:   "list <form-id>",
		Short: "List " + r.label,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			params := queryValues(cmd, paginationParams(cmd), "search")
			data, err := fetch(cmd, formPath(args[0], r.path), params)
			if err != nil {
				return err
			}
			printList(data, r.listKey, []output.Breadcrumb{
				{Label: "Show", Command: "neetoform forms " + r.use + " show <form-id> <" + r.singular + "-id>"},
			})
			return nil
		},
	}
	addPaginationFlags(list)
	if r.search {
		list.Flags().String("search", "", "Words to match in the rule name, description or condition field")
	}

	show := &cobra.Command{
		Use:   "show <form-id> <" + r.singular + "-id>",
		Short: "Show a " + r.singular,
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			data, err := fetch(cmd, formPath(args[0], r.path, args[1]), nil)
			if err != nil {
				return err
			}
			printResource(data, nil)
			return nil
		},
	}

	del := &cobra.Command{
		Use:   "delete <form-id> <" + r.singular + "-id>",
		Short: "Delete a " + r.singular,
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return removeAndConfirm(cmd, formPath(args[0], r.path, args[1]), url.Values{}, capitalize(r.singular)+" deleted.")
		},
	}

	clone := &cobra.Command{
		Use:   "clone <form-id> <" + r.singular + "-id>",
		Short: "Clone a " + r.singular,
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			data, err := send(cmd, http.MethodPost, formPath(args[0], r.path, args[1], "clone"), nil)
			if err != nil {
				return err
			}
			printActionResult(data, nil)
			return nil
		},
	}

	reorder := &cobra.Command{
		Use:   "reorder <form-id>",
		Short: "Set the order of " + r.label,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			ids, _ := cmd.Flags().GetString("ids")
			data, err := send(cmd, http.MethodPatch, formPath(args[0], r.path, "reorder"), map[string]any{r.idsKey: splitCSV(ids)})
			if err != nil {
				return err
			}
			printActionResult(data, nil)
			return nil
		},
	}
	reorder.Flags().String("ids", "", "Every "+r.singular+" id of the form in the new order, comma-separated")
	markFlagsRequired(reorder, "ids")

	group.AddCommand(list, show, del, clone, reorder)
	return group
}

func capitalize(s string) string {
	return string(s[0]-32) + s[1:]
}

func newLogicSave(r logicResource, verb string, build func(cmd *cobra.Command) (map[string]any, error)) *cobra.Command {
	args := cobra.ExactArgs(1)
	use := verb + " <form-id>"
	if verb == "update" {
		args = cobra.ExactArgs(2)
		use = verb + " <form-id> <" + r.singular + "-id>"
	}
	return &cobra.Command{
		Use:   use,
		Short: capitalize(verb) + " a " + r.singular,
		Args:  args,
		RunE: func(cmd *cobra.Command, args []string) error {
			body, err := build(cmd)
			if err != nil {
				return err
			}
			if verb == "update" {
				data, err := send(cmd, http.MethodPatch, formPath(args[0], r.path, args[1]), body)
				if err != nil {
					return err
				}
				printResource(data, nil)
				return nil
			}
			data, err := send(cmd, http.MethodPost, formPath(args[0], r.path), body)
			if err != nil {
				return err
			}
			printActionResult(data, nil)
			return nil
		},
	}
}

var conditionsResource = logicResource{
	use:      "conditions",
	path:     "conditions",
	singular: "condition",
	listKey:  "conditions",
	idsKey:   "condition_ids",
	label:    "conditional logic",
}

func init() {
	group := newLogicGroup(conditionsResource)
	for _, verb := range []string{"create", "update"} {
		save := newLogicSave(conditionsResource, verb, func(cmd *cobra.Command) (map[string]any, error) {
			return logicPayload(cmd, "condition", "enabled")
		})
		save.Flags().String("json-file", "", `Path to a JSON file with the condition, e.g. {"conditions":{"join_type":"and","items":[{"field":"<question-id>","verb":"is_equal","record":"custom_value","value":"Yes"}]},"actions":{"items":[{"action_type":"show","fields":["<question-id>"]}]}}`)
		save.Flags().Bool("enabled", true, "Turn the condition on or off")
		group.AddCommand(save)
	}
	formsCmd.AddCommand(group)
}
