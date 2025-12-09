#!/bin/bash
set -e

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
OUTPUT="$SCRIPT_DIR/coolpack"

# Version info
VERSION="${VERSION:-dev}"
COMMIT="$(git -C "$SCRIPT_DIR" rev-parse --short HEAD 2>/dev/null || echo "none")"
DATE="$(date -u +%Y-%m-%dT%H:%M:%SZ)"

# Build with ldflags
LDFLAGS="-s -w"
LDFLAGS="$LDFLAGS -X github.com/coollabsio/coolpack/pkg/version.Version=$VERSION"
LDFLAGS="$LDFLAGS -X github.com/coollabsio/coolpack/pkg/version.Commit=$COMMIT"
LDFLAGS="$LDFLAGS -X github.com/coollabsio/coolpack/pkg/version.Date=$DATE"

go build -ldflags "$LDFLAGS" -o "$OUTPUT" "$SCRIPT_DIR"

echo "$OUTPUT"
