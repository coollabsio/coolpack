package node

import (
	"regexp"
	"strings"

	"github.com/coollabsio/coolpack/pkg/app"
)

const DefaultNodeVersion = "24"

// DetectNodeVersion detects the Node.js version to use
// Priority:
// 1. COOLPACK_NODE_VERSION environment variable
// 2. NODE_VERSION environment variable
// 3. engines.node in package.json
// 4. .nvmrc file
// 5. .node-version file
// 6. .tool-versions file (asdf)
// 7. mise.toml file
// 8. Default to 22
func DetectNodeVersion(ctx *app.Context, pkg *PackageJSON) string {
	// 1. Check COOLPACK_NODE_VERSION env var
	if v := ctx.Env["COOLPACK_NODE_VERSION"]; v != "" {
		return normalizeVersion(v)
	}

	// 2. Check NODE_VERSION env var
	if v := ctx.Env["NODE_VERSION"]; v != "" {
		return normalizeVersion(v)
	}

	// 3. Check engines.node in package.json
	if pkg != nil && pkg.Engines.Node != "" {
		if v := parseEngineVersion(pkg.Engines.Node); v != "" {
			return v
		}
	}

	// 4. Check .nvmrc file
	if ctx.HasFile(".nvmrc") {
		if data, err := ctx.ReadFile(".nvmrc"); err == nil {
			if v := parseVersionFile(string(data)); v != "" {
				return v
			}
		}
	}

	// 5. Check .node-version file
	if ctx.HasFile(".node-version") {
		if data, err := ctx.ReadFile(".node-version"); err == nil {
			if v := parseVersionFile(string(data)); v != "" {
				return v
			}
		}
	}

	// 6. Check .tool-versions file (asdf format)
	if ctx.HasFile(".tool-versions") {
		if data, err := ctx.ReadFile(".tool-versions"); err == nil {
			if v := parseToolVersions(string(data), "nodejs"); v != "" {
				return v
			}
		}
	}

	// 7. Check mise.toml file
	if ctx.HasFile("mise.toml") {
		if data, err := ctx.ReadFile("mise.toml"); err == nil {
			if v := parseMiseToml(string(data)); v != "" {
				return v
			}
		}
	}

	// 8. Default
	return DefaultNodeVersion
}

// normalizeVersion cleans up version strings
func normalizeVersion(v string) string {
	v = strings.TrimSpace(v)
	v = strings.TrimPrefix(v, "v")
	return v
}

// parseVersionFile parses a simple version file (.nvmrc, .node-version)
func parseVersionFile(content string) string {
	v := strings.TrimSpace(content)
	v = strings.TrimPrefix(v, "v")

	// Handle lts/* or lts/iron type versions
	if strings.HasPrefix(strings.ToLower(v), "lts") {
		return DefaultNodeVersion
	}

	// Extract just the major version or full version
	if v != "" {
		return v
	}
	return ""
}

// parseEngineVersion parses a semver range from engines.node
// Examples: ">=18", "^20.0.0", "18.x", ">=18 <21"
func parseEngineVersion(constraint string) string {
	constraint = strings.TrimSpace(constraint)

	// Try to extract a version number
	re := regexp.MustCompile(`(\d+)(?:\.(\d+))?(?:\.(\d+))?`)
	matches := re.FindStringSubmatch(constraint)
	if len(matches) > 1 {
		// Return just the major version for broad compatibility
		return matches[1]
	}

	return ""
}

// parseToolVersions parses .tool-versions file (asdf format)
// Format: tool-name version
func parseToolVersions(content string, tool string) string {
	lines := strings.Split(content, "\n")
	for _, line := range lines {
		line = strings.TrimSpace(line)
		parts := strings.Fields(line)
		if len(parts) >= 2 && parts[0] == tool {
			return normalizeVersion(parts[1])
		}
	}
	return ""
}

// parseMiseToml extracts Node version from mise.toml
// Simple parser - just looks for node = "version"
func parseMiseToml(content string) string {
	// Look for patterns like:
	// node = "20"
	// nodejs = "20.10.0"
	// [tools]
	// node = "20"
	re := regexp.MustCompile(`(?:node|nodejs)\s*=\s*"([^"]+)"`)
	matches := re.FindStringSubmatch(content)
	if len(matches) > 1 {
		return normalizeVersion(matches[1])
	}
	return ""
}
