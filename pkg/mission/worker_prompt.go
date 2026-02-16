package mission

import (
	"fmt"
	"strings"
)

// BuildWorkerPrompt creates a minimal system prompt from a briefing.
// No personality, no skills, no memory — just precision execution instructions.
func BuildWorkerPrompt(b *Briefing) string {
	var sb strings.Builder

	sb.WriteString("You are a worker agent executing a specific task. Complete it precisely and efficiently.\n\n")

	sb.WriteString(fmt.Sprintf("## Task: %s\n\n", b.TaskID))
	sb.WriteString(fmt.Sprintf("**Objective:** %s\n\n", b.Objective))

	if b.Context != "" {
		sb.WriteString(fmt.Sprintf("**Context:** %s\n\n", b.Context))
	}

	if len(b.Constraints) > 0 {
		sb.WriteString("**Constraints:**\n")
		for _, c := range b.Constraints {
			sb.WriteString(fmt.Sprintf("- %s\n", c))
		}
		sb.WriteString("\n")
	}

	if len(b.AcceptanceCriteria) > 0 {
		sb.WriteString("**Acceptance Criteria:**\n")
		for _, ac := range b.AcceptanceCriteria {
			sb.WriteString(fmt.Sprintf("- %s\n", ac))
		}
		sb.WriteString("\n")
	}

	if len(b.FileScope) > 0 {
		sb.WriteString("**File Scope** (only modify these files):\n")
		for _, f := range b.FileScope {
			sb.WriteString(fmt.Sprintf("- `%s`\n", f))
		}
		sb.WriteString("\n")
	}

	if len(b.PredecessorFindings) > 0 {
		sb.WriteString("**Predecessor Findings:**\n")
		for _, pf := range b.PredecessorFindings {
			sb.WriteString(fmt.Sprintf("- %s\n", pf))
		}
		sb.WriteString("\n")
	}

	sb.WriteString("## Rules\n\n")
	sb.WriteString("- Complete the objective using the tools available to you.\n")
	sb.WriteString("- Only modify files within the specified scope (if given).\n")
	sb.WriteString("- When done, provide a concise summary of what was accomplished.\n")
	sb.WriteString("- If you encounter blockers, report them clearly.\n")

	return sb.String()
}

// BuildTaskMessage returns the initial user message for the worker loop.
func BuildTaskMessage(b *Briefing) string {
	return fmt.Sprintf("Execute task %s: %s", b.TaskID, b.Objective)
}
