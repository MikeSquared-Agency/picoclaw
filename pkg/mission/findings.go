package mission

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// Findings is the output report from a completed worker task.
type Findings struct {
	TaskID       string   `json:"task_id"`
	Summary      string   `json:"summary"`
	FilesChanged []string `json:"files_changed"`
	TestsRun     bool     `json:"tests_run"`
	TestsPassed  bool     `json:"tests_passed"`
	Issues       []string `json:"issues,omitempty"`
	NextSteps    []string `json:"next_steps,omitempty"`
}

// WriteFindings writes findings as markdown with JSON frontmatter to
// {missionDir}/.mission/findings/{taskID}.md.
func WriteFindings(missionDir string, f *Findings) error {
	dir := filepath.Join(missionDir, ".mission", "findings")
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("create findings dir: %w", err)
	}

	path := FindingsPath(missionDir, f.TaskID)

	frontmatter, err := json.MarshalIndent(f, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal findings: %w", err)
	}

	var md strings.Builder
	md.WriteString("---\n")
	md.Write(frontmatter)
	md.WriteString("\n---\n\n")
	md.WriteString("# Findings: " + f.TaskID + "\n\n")
	md.WriteString("## Summary\n\n" + f.Summary + "\n\n")

	if len(f.FilesChanged) > 0 {
		md.WriteString("## Files Changed\n\n")
		for _, file := range f.FilesChanged {
			md.WriteString("- `" + file + "`\n")
		}
		md.WriteString("\n")
	}

	if f.TestsRun {
		status := "PASSED"
		if !f.TestsPassed {
			status = "FAILED"
		}
		md.WriteString("## Tests\n\nTests were run: **" + status + "**\n\n")
	}

	if len(f.Issues) > 0 {
		md.WriteString("## Issues\n\n")
		for _, issue := range f.Issues {
			md.WriteString("- " + issue + "\n")
		}
		md.WriteString("\n")
	}

	if len(f.NextSteps) > 0 {
		md.WriteString("## Next Steps\n\n")
		for _, step := range f.NextSteps {
			md.WriteString("- " + step + "\n")
		}
		md.WriteString("\n")
	}

	if err := os.WriteFile(path, []byte(md.String()), 0644); err != nil {
		return fmt.Errorf("write findings: %w", err)
	}

	return nil
}

// FindingsPath returns the filesystem path for a findings file.
func FindingsPath(missionDir, taskID string) string {
	return filepath.Join(missionDir, ".mission", "findings", taskID+".md")
}
