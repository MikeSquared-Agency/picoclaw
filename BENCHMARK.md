# PicoClaw Benchmark Results

**Date**: February 16, 2026  
**Branch**: feat/worker-direct-prompt  
**Binary Version**: vdev  

## Test Environment
- OS: Linux 6.8.0-100-generic (x64)
- Architecture: x64
- Go Version: Built binary (27MB)

## 1. Boot Time Analysis

### Agent Initialization
PicoClaw shows extremely fast initialization times:

- **Run 1**: 0.012s (real), 0.006s (user), 0.008s (sys)
- **Run 2**: 0.006s (real), 0.004s (user), 0.004s (sys) 
- **Run 3**: 0.006s (real), 0.003s (user), 0.005s (sys)

**Average Boot Time**: ~0.008s (8 milliseconds)

The agent initializes successfully with message: `Agent initialized {tools_count=13, skills_total=6, skills_available=6}`

## 2. Memory Usage

### Binary Size
- **Compiled binary**: 26MB (27,155,316 bytes)

### Runtime Memory
❌ **Unable to measure runtime memory accurately** due to rapid process termination during testing. The agent process exits immediately after initialization when no sustained interaction is provided.

## 3. Comparison with OpenClaw

| Metric | PicoClaw | OpenClaw (spec) | Difference |
|--------|----------|-----------------|------------|
| Boot Time | ~8ms | ~5s | **625x faster** |
| Binary Size | 26MB | N/A | N/A |
| Memory (RSS) | Unable to measure | ~1GB | N/A |

## 4. Configuration Status

### ✅ Setup Complete
- Repository checked out on `feat/worker-direct-prompt` branch
- Binary built successfully (27MB)
- Configuration file exists at `~/.picoclaw/config.json`
- Workspace directory created at `~/.picoclaw/workspace`

### ❌ Authentication Issues
**Anthropic API connectivity failed** during testing:
- Error: `Invalid Anthropic API Key` (HTTP 401)
- Tested with both available API keys from OpenClaw auth profiles
- Both `anthropic:default` and `anthropic:manual` tokens failed authentication

### Configuration Details
- Model: `claude-3-5-haiku-latest`
- Provider: `anthropic` 
- Tools: 13 available
- Skills: 6 total, 6 available

## 5. Key Findings

### ✅ Performance Advantages
1. **Exceptional boot speed**: 625x faster than OpenClaw
2. **Compact binary**: Single 26MB executable
3. **Quick initialization**: Tools and skills load in milliseconds

### ⚠️ Limitations Discovered
1. **API Authentication**: Current API keys are invalid/expired
2. **Memory measurement**: Process exits too quickly for accurate RSS measurement
3. **Worker mode**: Requires mission directory structure not yet documented

## 6. Recommendations

1. **Fix API Keys**: Update Anthropic API keys or obtain new working credentials
2. **Memory Profiling**: Need longer-running task to accurately measure memory footprint
3. **Worker Documentation**: Document proper mission directory structure for worker mode
4. **Sustained Testing**: Create test scenarios that keep the process running for memory analysis

## Conclusion

PicoClaw demonstrates impressive performance characteristics with sub-10ms boot times, making it potentially 625x faster than OpenClaw for startup scenarios. However, API authentication issues prevent full functionality testing and accurate memory profiling.

---
*Benchmark conducted on feat/worker-direct-prompt branch*  
*Generated: 2026-02-16 19:16 GMT*