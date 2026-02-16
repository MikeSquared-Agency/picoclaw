package providers

import "testing"

func TestUsageInfoAdd(t *testing.T) {
	u := &UsageInfo{
		PromptTokens:     100,
		CompletionTokens: 50,
		TotalTokens:      150,
		CacheReadTokens:  20,
		CacheWriteTokens: 10,
	}

	other := &UsageInfo{
		PromptTokens:     200,
		CompletionTokens: 100,
		TotalTokens:      300,
		CacheReadTokens:  40,
		CacheWriteTokens: 20,
	}

	u.Add(other)

	if u.PromptTokens != 300 {
		t.Errorf("PromptTokens = %d, want 300", u.PromptTokens)
	}
	if u.CompletionTokens != 150 {
		t.Errorf("CompletionTokens = %d, want 150", u.CompletionTokens)
	}
	if u.TotalTokens != 450 {
		t.Errorf("TotalTokens = %d, want 450", u.TotalTokens)
	}
	if u.CacheReadTokens != 60 {
		t.Errorf("CacheReadTokens = %d, want 60", u.CacheReadTokens)
	}
	if u.CacheWriteTokens != 30 {
		t.Errorf("CacheWriteTokens = %d, want 30", u.CacheWriteTokens)
	}
}

func TestUsageInfoAddNil(t *testing.T) {
	u := &UsageInfo{PromptTokens: 100}
	u.Add(nil) // should not panic
	if u.PromptTokens != 100 {
		t.Errorf("PromptTokens = %d, want 100 (unchanged)", u.PromptTokens)
	}
}

func TestUsageInfoAddZero(t *testing.T) {
	u := &UsageInfo{}
	other := &UsageInfo{
		PromptTokens:     500,
		CompletionTokens: 200,
		TotalTokens:      700,
	}
	u.Add(other)
	if u.PromptTokens != 500 {
		t.Errorf("PromptTokens = %d, want 500", u.PromptTokens)
	}
	if u.TotalTokens != 700 {
		t.Errorf("TotalTokens = %d, want 700", u.TotalTokens)
	}
}
