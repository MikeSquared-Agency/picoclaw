package hermes

import (
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

const (
	SubjectSessionCompleted = "swarm.cc.session.completed"
	SubjectSessionFailed    = "swarm.cc.session.failed"
)

// Event is the standardised Hermes envelope, matching cc-sidecar and Warren.
type Event struct {
	ID        string          `json:"id"`
	Type      string          `json:"type"`
	Source    string          `json:"source"`
	Timestamp time.Time       `json:"timestamp"`
	Data      json.RawMessage `json:"data"`
}

// SessionData is the payload for cc.session.completed/failed events.
// Extends cc-sidecar's SessionData with token usage and runtime fields.
type SessionData struct {
	SessionID      string   `json:"session_id"`
	TaskID         string   `json:"task_id,omitempty"`
	OwnerUUID      string   `json:"owner_uuid,omitempty"`
	AgentType      string   `json:"agent_type"`
	TranscriptPath string   `json:"transcript_path,omitempty"`
	FilesChanged   []string `json:"files_changed"`
	ExitCode       int      `json:"exit_code"`
	DurationMs     int64    `json:"duration_ms"`
	WorkingDir     string   `json:"working_dir"`
	Timestamp      string   `json:"timestamp"`
	Model          string   `json:"model,omitempty"`
	Runtime        string   `json:"runtime,omitempty"`
	InputTokens    int64    `json:"input_tokens,omitempty"`
	OutputTokens   int64    `json:"output_tokens,omitempty"`
	CacheRead      int64    `json:"cache_read_tokens,omitempty"`
	CacheWrite     int64    `json:"cache_write_tokens,omitempty"`
}

// Client is a short-lived NATS JetStream publisher for PicoClaw worker events.
type Client struct {
	nc     *nats.Conn
	js     jetstream.JetStream
	logger *slog.Logger
}

// Connect creates a new Hermes client connected to NATS.
// Uses short-lived settings appropriate for worker processes (minutes, not hours).
func Connect(url, token string, logger *slog.Logger) (*Client, error) {
	opts := []nats.Option{
		nats.Name("picoclaw-worker"),
		nats.Timeout(5 * time.Second),
		nats.ReconnectWait(2 * time.Second),
		nats.MaxReconnects(3), // short-lived: fail fast
	}
	if token != "" {
		opts = append(opts, nats.Token(token))
	}

	nc, err := nats.Connect(url, opts...)
	if err != nil {
		return nil, fmt.Errorf("nats connect: %w", err)
	}

	js, err := jetstream.New(nc)
	if err != nil {
		nc.Close()
		return nil, fmt.Errorf("jetstream: %w", err)
	}

	return &Client{
		nc:     nc,
		js:     js,
		logger: logger,
	}, nil
}

// PublishCompleted publishes a session completed event.
func (c *Client) PublishCompleted(data *SessionData) error {
	return c.publish(SubjectSessionCompleted, "cc.session.completed", data)
}

// PublishFailed publishes a session failed event.
func (c *Client) PublishFailed(data *SessionData) error {
	return c.publish(SubjectSessionFailed, "cc.session.failed", data)
}

func (c *Client) publish(subject, eventType string, data *SessionData) error {
	raw, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("marshal session data: %w", err)
	}

	ev := Event{
		ID:        uuid.New().String(),
		Type:      eventType,
		Source:    "picoclaw-worker",
		Timestamp: time.Now().UTC(),
		Data:      raw,
	}

	evBytes, err := json.Marshal(ev)
	if err != nil {
		return fmt.Errorf("marshal event: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	ack, err := c.js.Publish(ctx, subject, evBytes)
	if err != nil {
		return fmt.Errorf("jetstream publish: %w", err)
	}

	c.logger.Info("published session event",
		"subject", subject,
		"task_id", data.TaskID,
		"stream", ack.Stream,
		"seq", ack.Sequence,
	)
	return nil
}

// Close drains and closes the NATS connection.
func (c *Client) Close() {
	if c.nc != nil {
		_ = c.nc.Drain()
	}
}
