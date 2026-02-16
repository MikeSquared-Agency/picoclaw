#!/bin/bash

echo "=== MEMORY USAGE BENCHMARK ==="

# Start picoclaw agent in background and capture its PID
echo "Starting picoclaw agent..."
cd ~/picoclaw
./picoclaw agent > /dev/null 2>&1 &
PICOCLAW_PID=$!

# Wait a moment for it to fully initialize
sleep 1

# Check if process is still running
if kill -0 $PICOCLAW_PID 2>/dev/null; then
    echo "PicoClaw process running with PID: $PICOCLAW_PID"
    
    # Get memory stats
    if [ -f "/proc/$PICOCLAW_PID/status" ]; then
        echo "Memory statistics:"
        grep -E "VmRSS|VmSize" /proc/$PICOCLAW_PID/status
        
        # Also get from ps
        echo ""
        echo "Memory from ps aux:"
        ps aux | grep $PICOCLAW_PID | grep -v grep
    else
        echo "Process status file not found"
    fi
    
    # Kill the process
    kill $PICOCLAW_PID 2>/dev/null
else
    echo "Process exited quickly"
fi