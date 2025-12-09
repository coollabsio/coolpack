package node

import (
	"encoding/json"
	"strings"
)

// PackageJSON represents the structure of a package.json file
type PackageJSON struct {
	Name             string            `json:"name"`
	Version          string            `json:"version"`
	Main             string            `json:"main"`
	Type             string            `json:"type"`
	Scripts          map[string]string `json:"scripts"`
	Dependencies     map[string]string `json:"dependencies"`
	DevDependencies  map[string]string `json:"devDependencies"`
	Engines          Engines           `json:"engines"`
	PackageManager   string            `json:"packageManager"`
	Workspaces       Workspaces        `json:"workspaces"`
	CacheDirectories []string          `json:"cacheDirectories"`
}

// Engines represents the engines field in package.json
type Engines struct {
	Node string `json:"node"`
	NPM  string `json:"npm"`
	Yarn string `json:"yarn"`
	PNPM string `json:"pnpm"`
	Bun  string `json:"bun"`
}

// Workspaces can be either an array of strings or an object with a packages field
type Workspaces struct {
	Packages []string
}

// UnmarshalJSON handles both array and object formats for workspaces
func (w *Workspaces) UnmarshalJSON(data []byte) error {
	// Try array first
	var arr []string
	if err := json.Unmarshal(data, &arr); err == nil {
		w.Packages = arr
		return nil
	}

	// Try object format
	var obj struct {
		Packages []string `json:"packages"`
	}
	if err := json.Unmarshal(data, &obj); err == nil {
		w.Packages = obj.Packages
		return nil
	}

	// Default to empty
	w.Packages = nil
	return nil
}

// ParsePackageJSON parses a package.json file from bytes
func ParsePackageJSON(data []byte) (*PackageJSON, error) {
	var pkg PackageJSON
	if err := json.Unmarshal(data, &pkg); err != nil {
		return nil, err
	}
	return &pkg, nil
}

// HasScript checks if a script exists in package.json
func (p *PackageJSON) HasScript(name string) bool {
	if p.Scripts == nil {
		return false
	}
	_, ok := p.Scripts[name]
	return ok
}

// GetScript returns a script from package.json
func (p *PackageJSON) GetScript(name string) string {
	if p.Scripts == nil {
		return ""
	}
	return p.Scripts[name]
}

// HasDependency checks if a dependency exists (in either dependencies or devDependencies)
func (p *PackageJSON) HasDependency(name string) bool {
	if p.Dependencies != nil {
		if _, ok := p.Dependencies[name]; ok {
			return true
		}
	}
	if p.DevDependencies != nil {
		if _, ok := p.DevDependencies[name]; ok {
			return true
		}
	}
	return false
}

// GetDependencyVersion returns the version of a dependency
func (p *PackageJSON) GetDependencyVersion(name string) string {
	if p.Dependencies != nil {
		if v, ok := p.Dependencies[name]; ok {
			return v
		}
	}
	if p.DevDependencies != nil {
		if v, ok := p.DevDependencies[name]; ok {
			return v
		}
	}
	return ""
}

// GetPackageManagerInfo parses the packageManager field (e.g., "pnpm@8.0.0")
// Returns the package manager name and version
func (p *PackageJSON) GetPackageManagerInfo() (name, version string) {
	if p.PackageManager == "" {
		return "", ""
	}

	parts := strings.SplitN(p.PackageManager, "@", 2)
	name = parts[0]
	if len(parts) > 1 {
		// Remove any hash suffix (e.g., "pnpm@8.0.0+sha256.xxx")
		version = strings.Split(parts[1], "+")[0]
	}
	return name, version
}

// IsMonorepo checks if this is a monorepo setup
func (p *PackageJSON) IsMonorepo() bool {
	return len(p.Workspaces.Packages) > 0
}
