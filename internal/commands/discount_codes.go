package commands

import (
	"net/http"
	"net/url"

	"github.com/neetozone/neeto-cli-commons/output"
	"github.com/spf13/cobra"
)

var validDiscountKinds = []string{"percentage", "fixed"}

var discountCodeBodySpecs = []string{"code", "value", "kind", "currency", "expires-at", "redeems-limit", "form-ids[]"}

var discountCodesCmd = &cobra.Command{
	Use:   "discount-codes",
	Short: "Manage discount codes",
}

func newDiscountCodeSaveCmd(use, short string, nargs int, method string, render printer) *cobra.Command {
	save := &cobra.Command{
		Use:   use,
		Short: short,
		Args:  exactArgs(nargs),
		RunE: func(cmd *cobra.Command, args []string) error {
			if err := validateEnum("kind", flagString(cmd, "kind"), validDiscountKinds); err != nil {
				return err
			}
			body := map[string]any{"discount_code": flagValues(cmd, discountCodeBodySpecs...)}
			return sendAndPrint(cmd, method, workspaceResource("discount_codes")(args[:nargs]), body, render)
		},
	}
	save.Flags().String("code", "", "Code customers enter, 3 to 30 characters without spaces")
	save.Flags().String("value", "", "Discount value; a percentage up to 100 or a fixed amount")
	save.Flags().String("kind", "", "Discount kind: percentage (default) or fixed")
	save.Flags().String("currency", "", "Currency of a fixed discount, e.g. usd")
	save.Flags().String("expires-at", "", "Expiry date in ISO 8601, e.g. 2026-12-31")
	save.Flags().Int("redeems-limit", 0, "Maximum number of redemptions, 1 to 100000")
	save.Flags().String("form-ids", "", "Comma-separated form IDs to enable the code on; on update, replaces the current forms")
	return save
}

var discountCodesReportCmd = &cobra.Command{
	Use:   "report <discount-code-id>",
	Short: "Show the redemption report of a discount code",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		params := queryValues(cmd, url.Values{}, "start-date", "end-date")
		data, err := fetch(cmd, apiPath("discount_codes", args[0], "report"), params)
		if err != nil {
			return err
		}
		printResource(data, nil)
		return nil
	},
}

func init() {
	createCmd := newDiscountCodeSaveCmd("create", "Create a discount code", 0, http.MethodPost, printActionResult)
	markFlagsRequired(createCmd, "code", "value")

	discountCodesFormsCmd := &cobra.Command{
		Use:   "forms",
		Short: "Work with forms a discount code can apply to",
	}
	discountCodesFormsCmd.AddCommand(listCmd("list", "List forms a discount code can be enabled on", 0,
		func([]string) string { return apiPath("discount_codes", "forms") }, "forms", nil))

	discountCodesRedemptionsCmd := &cobra.Command{
		Use:   "redemptions",
		Short: "Work with discount code redemptions",
	}
	discountCodesRedemptionsCmd.AddCommand(listCmd("list <discount-code-id>", "List redemptions of a discount code", 1,
		func(args []string) string { return apiPath("discount_codes", args[0], "redemptions") }, "redemptions", nil))

	discountCodesReportCmd.Flags().String("start-date", "", "First day to include, ISO date, e.g. 2026-10-01")
	discountCodesReportCmd.Flags().String("end-date", "", "Last day to include, ISO date, e.g. 2026-10-31")

	discountCodesCmd.AddCommand(
		listCmd("list", "List discount codes", 0, workspaceResource("discount_codes"), "discount_codes", []output.Breadcrumb{
			{Label: "Show", Command: "neetoform discount-codes show <id>"},
		}),
		showCmd("show <id>", "Show a discount code", 1, workspaceResource("discount_codes")),
		createCmd,
		newDiscountCodeSaveCmd("update <id>", "Update a discount code", 1, http.MethodPatch, printResource),
		deleteCmd("delete <id>", "Delete a discount code", 1, workspaceResource("discount_codes"), "Discount code deleted."),
		discountCodesFormsCmd,
		discountCodesRedemptionsCmd,
		discountCodesReportCmd,
	)
	register(func(root *cobra.Command) { root.AddCommand(discountCodesCmd) })
}
