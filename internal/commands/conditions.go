package commands

import (
	"net/http"
	"strings"

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

func newLogicGroup(r logicResource) *cobra.Command {
	group := &cobra.Command{Use: r.use, Short: "Manage " + r.label}
	idUse := " <form-id> <" + r.singular + "-id>"

	list := listCmd("list <form-id>", "List "+r.label, 1, formResource(r.path), r.listKey, []output.Breadcrumb{
		{Label: "Show", Command: "neetoform forms " + r.use + " show" + idUse},
	}, "search")
	if r.search {
		list.Flags().String("search", "", "Words to match in the rule name, description or condition field")
	}

	clone := &cobra.Command{
		Use:   "clone" + idUse,
		Short: "Clone a " + r.singular,
		Args:  cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			return sendAndPrint(cmd, http.MethodPost, formPath(args[0], r.path, args[1], "clone"), nil, printActionResult)
		},
	}

	reorder := &cobra.Command{
		Use:   "reorder <form-id>",
		Short: "Set the order of " + r.label,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			body := map[string]any{r.idsKey: splitCSV(flagString(cmd, "ids"))}
			return sendAndPrint(cmd, http.MethodPatch, formPath(args[0], r.path, "reorder"), body, printActionResult)
		},
	}
	reorder.Flags().String("ids", "", "Every "+r.singular+" id of the form in the new order, comma-separated")
	markFlagsRequired(reorder, "ids")

	group.AddCommand(
		list,
		showCmd("show"+idUse, "Show a "+r.singular, 2, formResource(r.path)),
		deleteCmd("delete"+idUse, "Delete a "+r.singular, 2, formResource(r.path), capitalize(r.singular)+" deleted."),
		clone,
		reorder,
	)
	return group
}

func capitalize(s string) string {
	return strings.ToUpper(s[:1]) + s[1:]
}

func newLogicSave(r logicResource, verb string, build func(cmd *cobra.Command) (map[string]any, error)) *cobra.Command {
	nargs, use, method, render := 1, verb+" <form-id>", http.MethodPost, printer(printActionResult)
	if verb == "update" {
		nargs, use, method, render = 2, verb+" <form-id> <"+r.singular+"-id>", http.MethodPatch, printResource
	}
	return &cobra.Command{
		Use:   use,
		Short: capitalize(verb) + " a " + r.singular,
		Args:  cobra.ExactArgs(nargs),
		RunE: func(cmd *cobra.Command, args []string) error {
			body, err := build(cmd)
			if err != nil {
				return err
			}
			return sendAndPrint(cmd, method, formResource(r.path)(args), body, render)
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
			condition, err := jsonFileWithFlags(cmd, "enabled")
			if err != nil {
				return nil, err
			}
			return map[string]any{"condition": condition}, nil
		})
		addJSONFileFlag(save, `Path to a JSON file with the condition, e.g. {"conditions":{"join_type":"and","items":[{"field":"<question-id>","verb":"is_equal","record":"custom_value","value":"Yes"}]},"actions":{"items":[{"action_type":"show","fields":["<question-id>"]}]}}`)
		save.Flags().Bool("enabled", true, "Turn the condition on or off")
		group.AddCommand(save)
	}
	formsCmd.AddCommand(group)
}
