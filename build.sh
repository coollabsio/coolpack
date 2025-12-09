#!/bin/bash
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
OUTPUT="$SCRIPT_DIR/coolpack"

go build -o "$OUTPUT" "$SCRIPT_DIR"

echo "$OUTPUT"
