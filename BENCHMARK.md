# PicoClaw Benchmark Report

**Date:** February 16, 2026  
**Branch:** `feat/worker-direct-prompt`  
**System:** Linux 6.8.0-100-generic (x64)  

## Summary

PicoClaw onboarding was partially successful, but worker functionality could not be fully tested due to API key authentication issues.

## Task 1: Onboarding Results

### ✅ Binary Status
- **Exists:** Yes, pre-built binary found at `~/picoclaw/picoclaw` (27.2MB)
- **Commands Available:** onboard, agent, worker, auth, gateway, status, cron, migrate, skills, version

### ✅ Configuration
- **Config File:** Already exists at `~/.picoclaw/config.json`
- **Provider:** Anthropic configured
- **Model:** `claude-3-5-haiku-latest` 
- **API Key:** Present in config (from OpenClaw auth-profiles.json)

### ❌ Authentication Issue
**Problem:** Anthropic API key validation failed
```
Error: LLM call failed: API request failed:
  Status: 401
  Body: {"error":{"code":"authentication_error","message":"Invalid Anthropic API Key","type":"invalid_request_error","param":null}}
```

**Root Cause:** The API key copied from OpenClaw's auth-profiles.json (`anthropic:manual` profile) is being rejected by Anthropic's API when called from PicoClaw.

**Possible Solutions:**
1. Generate a fresh Anthropic API key 
2. Check if key format/encoding differs between OpenClaw and PicoClaw
3. Verify network/proxy differences between the two systems

## Task 2: Benchmark Results (Partial)

### Boot Time Analysis

**Test Attempted:**
```bash
time ./picoclaw worker --task test-ping --system-prompt /tmp/system_prompt.txt --task-message "ping" --model claude-3-5-haiku-latest --debug
```

**Results:**
- **Startup Time:** ~0.004-0.006 seconds (binary launch + argument parsing)
- **Failure Point:** Authentication before API call
- **Total Duration:** Not measurable due to auth failure

### Memory Usage Analysis

**Binary Size:** 27.2MB (static binary)
**Process Memory:** Could not measure due to early auth failure

### Comparison with OpenClaw

| Metric | PicoClaw | OpenClaw | Notes |
|--------|----------|----------|-------|
| **Boot Time** | ~0.005s* | ~5s | *PicoClaw fails at auth, not full boot |
| **Memory** | Unknown | ~1GB | Could not measure PicoClaw due to auth failure |
| **Binary Size** | 27.2MB | Variable | PicoClaw is a single static binary |

**Key Differences Observed:**
1. **Architecture:** PicoClaw appears to be a single static Go binary vs OpenClaw's Node.js runtime
2. **Startup Speed:** PicoClaw binary launches almost instantaneously (~5ms) but OpenClaw includes full agent initialization (~5s)
3. **Direct Prompt Mode:** PicoClaw supports `--system-prompt` and `--task-message` for direct worker tasks without mission files

## Blocking Issues

### 1. Authentication Configuration
- **Status:** 🚫 Blocked
- **Issue:** Invalid Anthropic API key prevents worker execution
- **Impact:** Cannot measure actual worker performance or memory usage during LLM calls
- **Estimated Fix Time:** 5-10 minutes with valid API key

### 2. Worker Command Interface
- **Status:** ⚠️  Partial
- **Issue:** Direct prompt mode requires both `--system-prompt` and `--task-message` but still expects `--mission-dir`
- **Workaround:** Use briefing-based mission structure instead of direct prompts
- **Impact:** Different interface than expected from task description

## Next Steps

1. **Immediate:** Obtain valid Anthropic API key for PicoClaw
2. **Testing:** Re-run benchmark with working authentication
3. **Measurement:** Compare actual boot time including LLM initialization
4. **Memory:** Measure RSS during actual worker task execution
5. **Performance:** Test multiple task types to get representative metrics

## Technical Notes

- PicoClaw binary appears well-optimized (5ms launch time)
- Configuration system is working correctly
- Error handling and logging are clear and helpful
- Direct prompt mode exists but requires mission directory structure

---

**Conclusion:** PicoClaw shows promise for ultra-fast worker startup compared to OpenClaw, but authentication issues prevent complete benchmarking. The binary architecture appears significantly more lightweight than OpenClaw's Node.js stack.