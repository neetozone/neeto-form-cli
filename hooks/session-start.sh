#!/bin/sh
# NeetoForm CLI — session-start hook for Claude Code
# Lightweight auth liveness check. Always exits 0 (informational).

if ! command -v neetoform >/dev/null 2>&1; then
  echo "NeetoForm CLI is not installed or not on PATH."
  exit 0
fi

if neetoform whoami >/dev/null 2>&1; then
  echo "NeetoForm plugin active."
else
  echo "NeetoForm CLI installed but not authenticated. Run 'neetoform login' to authenticate."
fi

exit 0
