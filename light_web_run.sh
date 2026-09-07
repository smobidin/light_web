#!/usr/bin/env bash
set -euo pipefail

# ── configuration ──────────────────────────────────────────────────
PROJECT_DIR="${LIGHT_WEB_DIR:-/home/stas/dev/light_web}"
DEFAULT_IP="0.0.0.0"
DEFAULT_PORT=8081
DEFAULT_DIR="."

# ── parse args ─────────────────────────────────────────────────────
ip="$DEFAULT_IP"
port="$DEFAULT_PORT"
dir="$DEFAULT_DIR"

while [[ $# -gt 0 ]]; do
    case "$1" in
        --ip|-i) ip="$2"; shift 2 ;;
        --port|-p) port="$2"; shift 2 ;;
        --directory|-d) dir="$2"; shift 2 ;;
        --help|-h)
            echo "Usage: $(basename "$0") [--ip ADDR] [--port PORT] [--directory DIR]"
            exit 0
            ;;
        *) echo "Unknown: $1"; exit 1 ;;
    esac
done

# ── resolve paths ──────────────────────────────────────────────────
SERVER="$PROJECT_DIR/ultra_simple.py"
VENV_PYTHON="$PROJECT_DIR/.venv/bin/python3"

if [[ ! -f "$SERVER" ]]; then
    echo "Error: server not found at $SERVER"
    echo "Set LIGHT_WEB_DIR env var or edit PROJECT_DIR in this script."
    exit 1
fi

if [[ ! -x "$VENV_PYTHON" ]]; then
    echo "Error: venv python not found at $VENV_PYTHON"
    echo "Run: cd $PROJECT_DIR && python3 -m venv .venv && .venv/bin/pip install -r requirements.txt"
    exit 1
fi

# ── find accessible IPs ────────────────────────────────────────────
get_ips() {
    if command -v hostname &>/dev/null; then
        hostname -I 2>/dev/null
    elif command -v ip &>/dev/null; then
        ip -4 -o addr show scope global 2>/dev/null \
            | awk '{print $4}' | cut -d/ -f1
    fi
}

# ── run ────────────────────────────────────────────────────────────
echo "Starting Light Web Server..."
"$VENV_PYTHON" "$SERVER" --host "$ip" --port "$port" --directory "$dir" &
SRV_PID=$!

# give it a moment to bind
sleep 1

if ! kill -0 "$SRV_PID" 2>/dev/null; then
    echo "Server failed to start."
    exit 1
fi

echo ""
echo "  Local:    http://127.0.0.1:$port"
for addr in $(get_ips); do
    [[ -z "$addr" ]] && continue
    echo "  Network:  http://$addr:$port"
done
echo ""

trap 'echo ""; echo "Stopping server..."; kill "$SRV_PID" 2>/dev/null; wait "$SRV_PID" 2>/dev/null; exit 0' INT TERM

wait "$SRV_PID"
