package commands

import (
	"net/http"

	"github.com/neetozone/neeto-cli-commons/output"
	"github.com/spf13/cobra"
)

var (
	validTrackingMethods = []string{"no_track", "cookie_track", "ip_track", "record_track"}
	validAvailability    = []string{"enable", "disable", "based_on_submission_count", "based_on_time", "based_on_limit_and_time"}
	validAccessAuths     = []string{"no", "basic", "email", "session", "invited"}
)

type settingGroup struct {
	use       string
	short     string
	segment   string
	root      string
	configure func(update *cobra.Command)
	validate  func(cmd *cobra.Command) error
	specs     []string
}

func (g settingGroup) command() *cobra.Command {
	group := &cobra.Command{Use: g.use, Short: "Manage " + g.short}

	update := &cobra.Command{
		Use:   "update <form-id>",
		Short: "Update " + g.short,
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if g.validate != nil {
				if err := g.validate(cmd); err != nil {
					return err
				}
			}
			body := map[string]any{g.root: flagValues(cmd, g.specs...)}
			return sendAndPrint(cmd, http.MethodPatch, formPath(args[0], g.segment), body, printResource)
		},
	}
	g.configure(update)

	group.AddCommand(showCmd("show <form-id>", "Show "+g.short, 1, formResource(g.segment)), update)
	return group
}

var formDiscountCodesCmd = &cobra.Command{
	Use:   "discount-codes",
	Short: "Manage the discount codes a form accepts",
}

var formDiscountCodesUpdateCmd = &cobra.Command{
	Use:   "update <form-id>",
	Short: "Replace the discount codes a form accepts",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		body := map[string]any{"discount_codes": map[string]any{"ids": splitCSV(flagString(cmd, "ids"))}}
		return sendAndPrint(cmd, http.MethodPatch, formPath(args[0], "discount_codes"), body, printResource)
	},
}

func init() {
	groups := []settingGroup{
		{
			use:     "unique-submission",
			short:   "the unique submission setting",
			segment: "unique_submission",
			root:    "unique_submission",
			specs:   []string{"method=submission_tracking_method", "question-id=unique_by_record_id"},
			validate: func(cmd *cobra.Command) error {
				return validateEnum("method", flagString(cmd, "method"), validTrackingMethods)
			},
			configure: func(update *cobra.Command) {
				update.Flags().String("method", "", "How duplicate submissions are prevented: no_track, cookie_track, ip_track or record_track")
				update.Flags().String("question-id", "", "UUID of the email or phone question that identifies a respondent; used with record_track")
				markFlagsRequired(update, "method")
			},
		},
		{
			use:     "availability",
			short:   "when the form accepts submissions",
			segment: "availability",
			root:    "availability",
			specs:   []string{"method=manage_access_method", "deadline=submission_deadline", "limit=submission_limit", "closed-message=custom_closed_message"},
			validate: func(cmd *cobra.Command) error {
				return validateEnum("method", flagString(cmd, "method"), validAvailability)
			},
			configure: func(update *cobra.Command) {
				update.Flags().String("method", "", "enable, disable, based_on_submission_count, based_on_time or based_on_limit_and_time")
				update.Flags().String("deadline", "", "Stop accepting submissions at this ISO 8601 date-time, e.g. 2030-01-01T00:00:00Z")
				update.Flags().Int("limit", 0, "Stop accepting submissions after this many (1 to 1000000)")
				update.Flags().String("closed-message", "", "Message shown when the form is closed")
			},
		},
		{
			use:     "access-control",
			short:   "who can open the form",
			segment: "access_control",
			root:    "access_control",
			specs:   []string{"auth", "password", "allowed-keywords[]"},
			validate: func(cmd *cobra.Command) error {
				return validateEnum("auth", flagString(cmd, "auth"), validAccessAuths)
			},
			configure: func(update *cobra.Command) {
				update.Flags().String("auth", "", "no, basic (password), email (allowed email keywords), session (login) or invited")
				update.Flags().String("password", "", "Password for basic access, at least 4 characters")
				update.Flags().String("allowed-keywords", "", "Comma-separated allowed emails or domains for email access")
				markFlagsRequired(update, "auth")
			},
		},
		{
			use:     "tracking-codes",
			short:   "the tracking codes",
			segment: "tracking_codes",
			root:    "tracking_codes",
			specs:   []string{"google-analytics=ga_code", "google-tag-manager=gtm_code", "meta-pixel=pixel_id"},
			configure: func(update *cobra.Command) {
				update.Flags().String("google-analytics", "", "Google Analytics ID, G-... or UA-...; pass an empty value to remove")
				update.Flags().String("google-tag-manager", "", "Google Tag Manager ID, GTM-...; pass an empty value to remove")
				update.Flags().String("meta-pixel", "", "Meta Pixel ID, 15 or 16 digits; pass an empty value to remove")
			},
		},
		{
			use:     "neeto-chat",
			short:   "the NeetoChat widget setting",
			segment: "neeto_chat",
			root:    "neeto_chat",
			specs:   []string{"enabled"},
			configure: func(update *cobra.Command) {
				update.Flags().Bool("enabled", false, "Show the NeetoChat widget on the public form: --enabled=true or --enabled=false")
				markFlagsRequired(update, "enabled")
			},
		},
	}
	for _, group := range groups {
		formsCmd.AddCommand(group.command())
	}

	formDiscountCodesUpdateCmd.Flags().String("ids", "", "Comma-separated UUIDs of the discount codes the form accepts; omit to clear all")
	formDiscountCodesCmd.AddCommand(
		listCmd("list <form-id>", "List discount codes and whether the form accepts each", 1, formResource("discount_codes"), "discount_codes", []output.Breadcrumb{
			{Label: "Update", Command: "neetoform forms discount-codes update <form-id> --ids <code-id>"},
		}),
		formDiscountCodesUpdateCmd,
	)
	formsCmd.AddCommand(formDiscountCodesCmd)
}
