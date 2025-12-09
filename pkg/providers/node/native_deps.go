package node

// NativeDependency represents a Node.js package that requires native system dependencies
type NativeDependency struct {
	// Package is the npm package name
	Package string
	// AptPackages are the Debian/Ubuntu packages needed for building
	AptPackages []string
	// Description explains why these packages are needed
	Description string
}

// NativeDependencies is a list of known packages requiring native dependencies
var NativeDependencies = []NativeDependency{
	{
		Package:     "sharp",
		AptPackages: []string{"libvips-dev"},
		Description: "Image processing library",
	},
	{
		Package:     "@prisma/client",
		AptPackages: []string{"openssl"},
		Description: "Database ORM",
	},
	{
		Package:     "prisma",
		AptPackages: []string{"openssl"},
		Description: "Database ORM CLI",
	},
	{
		Package:     "puppeteer",
		AptPackages: []string{
			"chromium",
			"libnss3",
			"libatk1.0-0",
			"libatk-bridge2.0-0",
			"libcups2",
			"libdrm2",
			"libxkbcommon0",
			"libxcomposite1",
			"libxdamage1",
			"libxfixes3",
			"libxrandr2",
			"libgbm1",
			"libasound2",
			"libpango-1.0-0",
			"libcairo2",
		},
		Description: "Headless Chrome automation",
	},
	{
		Package:     "playwright",
		AptPackages: []string{
			"libnss3",
			"libatk1.0-0",
			"libatk-bridge2.0-0",
			"libcups2",
			"libdrm2",
			"libxkbcommon0",
			"libxcomposite1",
			"libxdamage1",
			"libxfixes3",
			"libxrandr2",
			"libgbm1",
			"libasound2",
			"libpango-1.0-0",
			"libcairo2",
		},
		Description: "Browser automation",
	},
	{
		Package:     "canvas",
		AptPackages: []string{"libcairo2-dev", "libjpeg-dev", "libpango1.0-dev", "libgif-dev", "librsvg2-dev"},
		Description: "Canvas rendering",
	},
	{
		Package:     "bcrypt",
		AptPackages: []string{"build-essential", "python3"},
		Description: "Password hashing",
	},
	{
		Package:     "argon2",
		AptPackages: []string{"build-essential"},
		Description: "Password hashing",
	},
	{
		Package:     "sqlite3",
		AptPackages: []string{"build-essential", "python3"},
		Description: "SQLite database",
	},
	{
		Package:     "better-sqlite3",
		AptPackages: []string{"build-essential", "python3"},
		Description: "SQLite database",
	},
	{
		Package:     "node-gyp",
		AptPackages: []string{"build-essential", "python3"},
		Description: "Native addon build tool",
	},
	{
		Package:     "cpu-features",
		AptPackages: []string{"build-essential"},
		Description: "CPU feature detection",
	},
	{
		Package:     "ssh2",
		AptPackages: []string{"build-essential"},
		Description: "SSH client",
	},
	{
		Package:     "libsql",
		AptPackages: []string{"build-essential"},
		Description: "LibSQL database",
	},
	{
		Package:     "@libsql/client",
		AptPackages: []string{"build-essential"},
		Description: "LibSQL client",
	},
}

// DetectNativeDependencies checks which native dependencies are used by the project
func DetectNativeDependencies(pkg *PackageJSON) []NativeDependency {
	var detected []NativeDependency

	for _, dep := range NativeDependencies {
		if pkg.HasDependency(dep.Package) {
			detected = append(detected, dep)
		}
	}

	return detected
}

// GetRequiredAptPackages returns a deduplicated list of APT packages needed
func GetRequiredAptPackages(deps []NativeDependency) []string {
	seen := make(map[string]bool)
	var packages []string

	for _, dep := range deps {
		for _, pkg := range dep.AptPackages {
			if !seen[pkg] {
				seen[pkg] = true
				packages = append(packages, pkg)
			}
		}
	}

	return packages
}
