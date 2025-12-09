package node

import (
	"github.com/coollabsio/coolpack/pkg/app"
)

// PackageManager represents a Node.js package manager
type PackageManager string

const (
	PackageManagerNPM       PackageManager = "npm"
	PackageManagerYarn1     PackageManager = "yarn"
	PackageManagerYarnBerry PackageManager = "yarnberry"
	PackageManagerPNPM      PackageManager = "pnpm"
	PackageManagerBun       PackageManager = "bun"
)

// PackageManagerInfo contains information about the detected package manager
type PackageManagerInfo struct {
	Name    PackageManager
	Version string
}

// DetectPackageManager detects the package manager used by the project
// Detection priority:
// 1. packageManager field in package.json
// 2. Lock files
// 3. engines field in package.json
// 4. Default to npm
func DetectPackageManager(ctx *app.Context, pkg *PackageJSON) PackageManagerInfo {
	info := PackageManagerInfo{
		Name:    PackageManagerNPM,
		Version: "",
	}

	// 1. Check packageManager field in package.json
	if pmName, pmVersion := pkg.GetPackageManagerInfo(); pmName != "" {
		switch pmName {
		case "pnpm":
			info.Name = PackageManagerPNPM
			info.Version = pmVersion
			return info
		case "yarn":
			// Check if it's Yarn Berry (2+)
			if isYarnBerry(pmVersion) {
				info.Name = PackageManagerYarnBerry
			} else {
				info.Name = PackageManagerYarn1
			}
			info.Version = pmVersion
			return info
		case "bun":
			info.Name = PackageManagerBun
			info.Version = pmVersion
			return info
		case "npm":
			info.Name = PackageManagerNPM
			info.Version = pmVersion
			return info
		}
	}

	// 2. Check lock files
	if ctx.HasFile("pnpm-lock.yaml") {
		info.Name = PackageManagerPNPM
		return info
	}

	if ctx.HasFile("bun.lockb") || ctx.HasFile("bun.lock") {
		info.Name = PackageManagerBun
		return info
	}

	// Check for Yarn Berry (.yarnrc.yml indicates Yarn 2+)
	if ctx.HasFile(".yarnrc.yml") || ctx.HasFile(".yarnrc.yaml") {
		info.Name = PackageManagerYarnBerry
		return info
	}

	if ctx.HasFile("yarn.lock") {
		info.Name = PackageManagerYarn1
		return info
	}

	if ctx.HasFile("package-lock.json") {
		info.Name = PackageManagerNPM
		return info
	}

	// 3. Check engines field
	if pkg.Engines.PNPM != "" {
		info.Name = PackageManagerPNPM
		return info
	}
	if pkg.Engines.Bun != "" {
		info.Name = PackageManagerBun
		return info
	}
	if pkg.Engines.Yarn != "" {
		info.Name = PackageManagerYarn1
		return info
	}

	// 4. Default to npm
	return info
}

// isYarnBerry checks if the version indicates Yarn 2+
func isYarnBerry(version string) bool {
	if version == "" {
		return false
	}
	// Yarn 2+ starts with 2., 3., 4., etc.
	if len(version) > 0 && version[0] >= '2' && version[0] <= '9' {
		return true
	}
	return false
}

// GetInstallCommand returns the install command for the package manager
func (pm PackageManagerInfo) GetInstallCommand() string {
	switch pm.Name {
	case PackageManagerPNPM:
		return "pnpm install --frozen-lockfile"
	case PackageManagerYarnBerry:
		return "yarn install --immutable"
	case PackageManagerYarn1:
		return "yarn install --frozen-lockfile"
	case PackageManagerBun:
		return "bun install --frozen-lockfile"
	default:
		return "npm ci"
	}
}

// GetRunCommand returns the run command prefix for the package manager
func (pm PackageManagerInfo) GetRunCommand() string {
	switch pm.Name {
	case PackageManagerPNPM:
		return "pnpm"
	case PackageManagerYarnBerry, PackageManagerYarn1:
		return "yarn"
	case PackageManagerBun:
		return "bun run"
	default:
		return "npm run"
	}
}

// GetLockFile returns the lock file name for the package manager
func (pm PackageManagerInfo) GetLockFile() string {
	switch pm.Name {
	case PackageManagerPNPM:
		return "pnpm-lock.yaml"
	case PackageManagerYarnBerry, PackageManagerYarn1:
		return "yarn.lock"
	case PackageManagerBun:
		return "bun.lockb"
	default:
		return "package-lock.json"
	}
}
