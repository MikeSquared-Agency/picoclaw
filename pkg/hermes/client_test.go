package hermes

import (
	"encoding/json"
	"testing"
)

func TestSessionDataJSON(t *testing.T) {
	data := SessionData{
		SessionID:    "picoclaw-test-1",
		TaskID:       "task-42",
		AgentType:    "picoclaw",
		FilesChanged: []string{"main.go", "config.go"},
		ExitCode:     0,
		DurationMs:   12345,
		WorkingDir:   "/tmp/test",
		Timestamp:    "2026-02-16T00:00:00Z",
		Model:        "anthropic/claude-sonnet-4",
		Runtime:      "picoclaw-dev",
		InputTokens:  1000,
		OutputTokens: 500,
		CacheRead:    200,
		CacheWrite:   100,
	}

	raw, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}

	// Verify it matches cc-sidecar format (base fields)
	var m map[string]interface{}
	if err := json.Unmarshal(raw, &m); err != nil {
		t.Fatal(err)
	}

	// Check base fields match cc-sidecar format
	if m["session_id"] != "picoclaw-test-1" {
		t.Errorf("session_id = %v", m["session_id"])
	}
	if m["agent_type"] != "picoclaw" {
		t.Errorf("agent_type = %v", m["agent_type"])
	}
	if m["task_id"] != "task-42" {
		t.Errorf("task_id = %v", m["task_id"])
	}

	// Check extended fields
	if m["model"] != "anthropic/claude-sonnet-4" {
		t.Errorf("model = %v", m["model"])
	}
	if m["runtime"] != "picoclaw-dev" {
		t.Errorf("runtime = %v", m["runtime"])
	}
	if m["input_tokens"] != float64(1000) {
		t.Errorf("input_tokens = %v", m["input_tokens"])
	}
	if m["output_tokens"] != float64(500) {
		t.Errorf("output_tokens = %v", m["output_tokens"])
	}
}

func TestSessionDataJSON_OmitsEmpty(t *testing.T) {
	// Minimal data — optional fields should be omitted
	data := SessionData{
		SessionID:    "picoclaw-test-2",
		AgentType:    "picoclaw",
		FilesChanged: []string{},
		Timestamp:    "2026-02-16T00:00:00Z",
	}

	raw, err := json.Marshal(data)
	if err != nil {
		t.Fatal(err)
	}

	var m map[string]interface{}
	json.Unmarshal(raw, &m)

	// These should be absent (omitempty)
	for _, key := range []string{"task_id", "owner_uuid", "transcript_path", "model", "runtime", "input_tokens", "output_tokens", "cache_read_tokens", "cache_write_tokens"} {
		if _, ok := m[key]; ok {
			t.Errorf("expected %s to be omitted, got %v", key, m[key])
		}
	}
}

func TestEventEnvelopeFormat(t *testing.T) {
	data := SessionData{
		SessionID: "test-session",
		AgentType: "picoclaw",
		Timestamp: "2026-02-16T00:00:00Z",
	}

	raw, _ := json.Marshal(data)

	ev := Event{
		ID:        "test-id",
		Type:      "cc.session.completed",
		Source:    "picoclaw-worker",
		Data:      raw,
	}

	evBytes, err := json.Marshal(ev)
	if err != nil {
		t.Fatal(err)
	}

	var m map[string]interface{}
	json.Unmarshal(evBytes, &m)

	if m["id"] != "test-id" {
		t.Errorf("id = %v", m["id"])
	}
	if m["type"] != "cc.session.completed" {
		t.Errorf("type = %v", m["type"])
	}
	if m["source"] != "picoclaw-worker" {
		t.Errorf("source = %v", m["source"])
	}

	// Data should be nested JSON
	if _, ok := m["data"]; !ok {
		t.Error("expected data field")
	}
}
