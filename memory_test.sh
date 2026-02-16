#!/bin/bash

echo "=== MEMORY USAGE TEST ==="

# Start the agent in the background
cd ~/picoclaw
(echo "test" | timeout 10s ./picoclaw agent) &
AGENT_PID=$!

# Wait for it to start
sleep 0.1

# Try to get memory info quickly
ps aux | grep picoclaw | grep -v grep | head -1

# Wait for it to finish
wait $AGENT_PID 2>/dev/null