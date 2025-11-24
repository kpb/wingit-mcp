package prompts

import "fmt"

// BuildFieldChecklistPrompt returns the model-facing text for the
// "field_checklist" prompt. Tolerates empty inputs by supplying
// friendly defaults so callers don't have to pre-validate.
func BuildFieldChecklistPrompt(location, dayRange string) string {
	if location == "" {
		location = "this area"
	}
	if dayRange == "" {
		dayRange = "the recent period"
	}

	return fmt.Sprintf(
		`You are a birding assistant. The user has just called the WingIt-MCP tool "target_checklist" to get likely new lifers near %s for %s.

Using the tool output provided in this conversation (JSON with "targets" and "filters"), produce a concise, printable field checklist:

- Focus only on likely lifers (the "targets" array).
- Group species by approximate recent frequency (high / medium / low) based on "recentFrequency".
- For each species, show: common name, scientific name, and a short note like "seen recently at <locName>" if present.
- Keep it compact, suitable for printing or quick reference in the field.
- Do not reprint the raw JSON; summarize it.

If there are no targets, explain that there are no likely new lifers for this query and suggest broadening radius or daysBack.`, location, dayRange)
}
