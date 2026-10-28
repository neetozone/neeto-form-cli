package commands

import (
	"github.com/spf13/cobra"
)

var (
	validRuleStatuses = []string{"enabled", "disabled"}
	validDelayUnits   = []string{"minutes", "hours", "days"}
	validRuleEvents   = []string{"submission_created", "custom_field_updated"}
)

var automationRulesResource = logicResource{
	use:      "automation-rules",
	path:     "automation_rules",
	singular: "rule",
	listKey:  "automation_rules",
	idsKey:   "rule_ids",
	label:    "automation rules",
	search:   true,
}

func buildRulePayload(cmd *cobra.Command) (map[string]any, error) {
	if err := validateEnum("status", flagString(cmd, "status"), validRuleStatuses); err != nil {
		return nil, err
	}
	if err := validateEnum("delay-unit", flagString(cmd, "delay-unit"), validDelayUnits); err != nil {
		return nil, err
	}
	if err := validateEnumList("events", splitCSV(flagString(cmd, "events")), validRuleEvents); err != nil {
		return nil, err
	}

	rule, err := jsonFileWithFlags(cmd, "name", "description", "status", "performer", "events[]")
	if err != nil {
		return nil, err
	}
	mergeInto(rule, "execution_delay", flagValues(cmd, "delay-value=value", "delay-unit=unit"))
	return map[string]any{"automation_rule": rule}, nil
}

func addRuleFlags(cmd *cobra.Command) {
	cmd.Flags().String("name", "", "Rule name, unique within the form")
	cmd.Flags().String("description", "", "Rule description")
	cmd.Flags().String("status", "", "enabled or disabled")
	cmd.Flags().String("performer", "", "any, or an organization role name that must perform the action")
	cmd.Flags().String("events", "", "Comma-separated events: submission_created, custom_field_updated")
	cmd.Flags().Int("delay-value", 0, "Delay before the actions run, 1 to 15")
	cmd.Flags().String("delay-unit", "", "Delay unit: minutes, hours or days")
	addJSONFileFlag(cmd, `Path to a JSON file with the rule; flags override it, e.g. {"conditions":[{"field":"<question-id>","kind":"record","verb":"contains","join_type":"and_operator","metadata":{"value":"@example.com"}}],"actions":[{"name":"email_to","metadata":{"subject":"New lead","body":"{{all-answers}}","emails":["sales@example.com"]}}]}`)
}

func init() {
	group := newLogicGroup(automationRulesResource)
	for _, verb := range []string{"create", "update"} {
		save := newLogicSave(automationRulesResource, verb, buildRulePayload)
		addRuleFlags(save)
		if verb == "create" {
			markFlagsRequired(save, "name")
		}
		group.AddCommand(save)
	}
	formsCmd.AddCommand(group)
}
