package mission

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// Briefing is the input specification for a worker task.
type Briefing struct {
	TaskID              string   `json:"task_id"`
	MissionID           string   `json:"mission_id,omitempty"`
	Objective           string   `json:"objective"`
	Context             string   `json:"context,omitempty"`
	Constraints         []string `json:"constraints,omitempty"`
	AcceptanceCriteria  []string `json:"acceptance_criteria,omitempty"`
	PredecessorFindings []string `json:"predecessor_findings,omitempty"`
	FileScope           []string `json:"file_scope,omitempty"`
}

// ReadBriefing loads a briefing from {missionDir}/.mission/handoffs/{taskID}-briefing.json.
func ReadBriefing(missionDir, taskID string) (*Briefing, error) {
	path := BriefingPath(missionDir, taskID)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read briefing %s: %w", path, err)
	}

	var b Briefing
	if err := json.Unmarshal(data, &b); err != nil {
		return nil, fmt.Errorf("parse briefing: %w", err)
	}

	if b.TaskID == "" {
		b.TaskID = taskID
	}

	return &b, nil
}

// BriefingPath returns the filesystem path for a briefing file.
func BriefingPath(missionDir, taskID string) string {
	return filepath.Join(missionDir, ".mission", "handoffs", taskID+"-briefing.json")
}
