package mission

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestReadBriefing(t *testing.T) {
	dir := t.TempDir()
	taskID := "test-task-1"

	// Create briefing directory
	handoffsDir := filepath.Join(dir, ".mission", "handoffs")
	if err := os.MkdirAll(handoffsDir, 0755); err != nil {
		t.Fatal(err)
	}

	briefing := Briefing{
		TaskID:             taskID,
		Objective:          "Create hello.txt with Hello World",
		Context:            "test context",
		Constraints:        []string{"no external deps"},
		AcceptanceCriteria: []string{"hello.txt exists"},
		FileScope:          []string{"hello.txt"},
	}

	data, err := json.Marshal(briefing)
	if err != nil {
		t.Fatal(err)
	}

	path := BriefingPath(dir, taskID)
	if err := os.WriteFile(path, data, 0644); err != nil {
		t.Fatal(err)
	}

	// Read it back
	got, err := ReadBriefing(dir, taskID)
	if err != nil {
		t.Fatal(err)
	}

	if got.TaskID != taskID {
		t.Errorf("TaskID = %q, want %q", got.TaskID, taskID)
	}
	if got.Objective != briefing.Objective {
		t.Errorf("Objective = %q, want %q", got.Objective, briefing.Objective)
	}
	if len(got.FileScope) != 1 || got.FileScope[0] != "hello.txt" {
		t.Errorf("FileScope = %v, want [hello.txt]", got.FileScope)
	}
}

func TestReadBriefing_FillsTaskID(t *testing.T) {
	dir := t.TempDir()
	handoffsDir := filepath.Join(dir, ".mission", "handoffs")
	os.MkdirAll(handoffsDir, 0755)

	// Briefing without task_id
	data := []byte(`{"objective":"do something"}`)
	path := BriefingPath(dir, "auto-id")
	os.WriteFile(path, data, 0644)

	got, err := ReadBriefing(dir, "auto-id")
	if err != nil {
		t.Fatal(err)
	}
	if got.TaskID != "auto-id" {
		t.Errorf("TaskID = %q, want %q", got.TaskID, "auto-id")
	}
}

func TestReadBriefing_NotFound(t *testing.T) {
	_, err := ReadBriefing("/nonexistent", "no-such-task")
	if err == nil {
		t.Error("expected error for missing briefing")
	}
}

func TestWriteFindings(t *testing.T) {
	dir := t.TempDir()

	findings := &Findings{
		TaskID:       "test-task-1",
		Summary:      "Created the file successfully",
		FilesChanged: []string{"hello.txt", "world.txt"},
		TestsRun:     true,
		TestsPassed:  true,
		Issues:       []string{"minor warning"},
		NextSteps:    []string{"review changes"},
	}

	if err := WriteFindings(dir, findings); err != nil {
		t.Fatal(err)
	}

	path := FindingsPath(dir, "test-task-1")
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	content := string(data)

	// Check frontmatter
	if !strings.HasPrefix(content, "---\n") {
		t.Error("expected frontmatter start")
	}
	if !strings.Contains(content, `"task_id": "test-task-1"`) {
		t.Error("expected task_id in frontmatter")
	}

	// Check markdown sections
	if !strings.Contains(content, "# Findings: test-task-1") {
		t.Error("expected findings header")
	}
	if !strings.Contains(content, "## Summary") {
		t.Error("expected summary section")
	}
	if !strings.Contains(content, "`hello.txt`") {
		t.Error("expected file listed")
	}
	if !strings.Contains(content, "**PASSED**") {
		t.Error("expected test status")
	}
}

func TestBuildWorkerPrompt(t *testing.T) {
	b := &Briefing{
		TaskID:             "task-42",
		Objective:          "Fix the bug",
		Context:            "There is a crash in main.go",
		Constraints:        []string{"don't modify tests"},
		AcceptanceCriteria: []string{"no crash on startup"},
		FileScope:          []string{"main.go"},
		PredecessorFindings: []string{"logs show nil pointer"},
	}

	prompt := BuildWorkerPrompt(b)

	if !strings.Contains(prompt, "task-42") {
		t.Error("expected task ID in prompt")
	}
	if !strings.Contains(prompt, "Fix the bug") {
		t.Error("expected objective in prompt")
	}
	if !strings.Contains(prompt, "don't modify tests") {
		t.Error("expected constraint in prompt")
	}
	if !strings.Contains(prompt, "`main.go`") {
		t.Error("expected file scope in prompt")
	}
	if !strings.Contains(prompt, "nil pointer") {
		t.Error("expected predecessor findings in prompt")
	}
}

func TestBuildTaskMessage(t *testing.T) {
	b := &Briefing{
		TaskID:    "task-42",
		Objective: "Fix the bug",
	}
	msg := BuildTaskMessage(b)
	if !strings.Contains(msg, "task-42") {
		t.Error("expected task ID in message")
	}
	if !strings.Contains(msg, "Fix the bug") {
		t.Error("expected objective in message")
	}
}
