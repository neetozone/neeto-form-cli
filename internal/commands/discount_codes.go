package commands

import (
	"net/url"

	"github.com/neetozone/neeto-cli-commons/output"
	"github.com/spf13/cobra"
)

var discountCodeBodySpecs = []string{"code", "value", "kind", "currency", "expires-at", "redeems-limit", "form-ids[]"}

var discountCodesCmd = &cobra.Command{
	Use:   "discount-codes",
	Short: "Manage discount codes",
}

var discountCodesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List discount codes",
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := fetch(cmd, apiPath("discount_codes"), paginationParams(cmd))
		if err != nil {
			return err
		}
		printList(data, "discount_codes", []output.Breadcrumb{
			{Label: "Show", Command: "neetoform discount-codes show <id>"},
		})
		return nil
	},
}

var discountCodesShowCmd = &cobra.Command{
	Use:   "show <id>",
	Short: "Show a discount code",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := fetch(cmd, apiPath("discount_codes", args[0]), nil)
		if err != nil {
			return err
		}
		printResource(data, nil)
		return nil
	},
}

var discountCodesCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a discount code",
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := validateEnum("kind", flagString(cmd, "kind"), validDiscountKinds); err != nil {
			return err
		}
		body := map[string]any{"discount_code": flagValues(cmd, discountCodeBodySpecs...)}
		data, err := send(cmd, "POST", apiPath("discount_codes"), body)
		if err != nil {
			return err
		}
		printActionResult(data, nil)
		return nil
	},
}

var discountCodesUpdateCmd = &cobra.Command{
	Use:   "update <id>",
	Short: "Update a discount code",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := validateEnum("kind", flagString(cmd, "kind"), validDiscountKinds); err != nil {
			return err
		}
		body := map[string]any{"discount_code": flagValues(cmd, discountCodeBodySpecs...)}
		data, err := send(cmd, "PATCH", apiPath("discount_codes", args[0]), body)
		if err != nil {
			return err
		}
		printResource(data, nil)
		return nil
	},
}

var discountCodesDeleteCmd = &cobra.Command{
	Use:   "delete <id>",
	Short: "Delete a discount code",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		return removeAndConfirm(cmd, apiPath("discount_codes", args[0]), nil, "Discount code deleted.")
	},
}

var discountCodesFormsCmd = &cobra.Command{
	Use:   "forms",
	Short: "Work with forms a discount code can apply to",
}

var discountCodesFormsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List forms a discount code can be enabled on",
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := fetch(cmd, apiPath("discount_codes", "forms"), paginationParams(cmd))
		if err != nil {
			return err
		}
		printList(data, "forms", nil)
		return nil
	},
}

var discountCodesRedemptionsCmd = &cobra.Command{
	Use:   "redemptions",
	Short: "Work with discount code redemptions",
}

var discountCodesRedemptionsListCmd = &cobra.Command{
	Use:   "list <discount-code-id>",
	Short: "List redemptions of a discount code",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		data, err := fetch(cmd, apiPath("discount_codes", args[0], "redemptions"), paginationParams(cmd))
		if err != nil {
			return err
		}
		printList(data, "redemptions", nil)
		return nil
	},
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

var validDiscountKinds = []string{"percentage", "fixed"}

func addDiscountCodeFlags(cmd *cobra.Command) {
	cmd.Flags().String("code", "", "Code customers enter, 3 to 30 characters without spaces")
	cmd.Flags().String("value", "", "Discount value; a percentage up to 100 or a fixed amount")
	cmd.Flags().String("kind", "", "Discount kind: percentage (default) or fixed")
	cmd.Flags().String("currency", "", "Currency of a fixed discount, e.g. usd")
	cmd.Flags().String("expires-at", "", "Expiry date in ISO 8601, e.g. 2026-12-31")
	cmd.Flags().Int("redeems-limit", 0, "Maximum number of redemptions, 1 to 100000")
	cmd.Flags().String("form-ids", "", "Comma-separated form IDs to enable the code on; on update, replaces the current forms")
}

func init() {
	addPaginationFlags(discountCodesListCmd)
	addPaginationFlags(discountCodesFormsListCmd)
	addPaginationFlags(discountCodesRedemptionsListCmd)

	addDiscountCodeFlags(discountCodesCreateCmd)
	markFlagsRequired(discountCodesCreateCmd, "code", "value")
	addDiscountCodeFlags(discountCodesUpdateCmd)

	discountCodesReportCmd.Flags().String("start-date", "", "First day to include, ISO date, e.g. 2026-10-01")
	discountCodesReportCmd.Flags().String("end-date", "", "Last day to include, ISO date, e.g. 2026-10-31")

	discountCodesFormsCmd.AddCommand(discountCodesFormsListCmd)
	discountCodesRedemptionsCmd.AddCommand(discountCodesRedemptionsListCmd)
	discountCodesCmd.AddCommand(discountCodesListCmd, discountCodesShowCmd, discountCodesCreateCmd,
		discountCodesUpdateCmd, discountCodesDeleteCmd, discountCodesFormsCmd,
		discountCodesRedemptionsCmd, discountCodesReportCmd)
	register(func(root *cobra.Command) { root.AddCommand(discountCodesCmd) })
}
