package node

import (
	"fmt"

	"github.com/coollabsio/coolpack/pkg/app"
)

// Provider is the Node.js provider implementation
type Provider struct{}

// New creates a new Node.js provider
func New() *Provider {
	return &Provider{}
}

// Name returns the provider name
func (p *Provider) Name() string {
	return "node"
}

// Detect checks if the application is a Node.js project
func (p *Provider) Detect(ctx *app.Context) (bool, error) {
	return ctx.HasFile("package.json"), nil
}

// Plan generates a build plan for the Node.js application
func (p *Provider) Plan(ctx *app.Context) (*app.Plan, error) {
	// Read and parse package.json
	pkgData, err := ctx.ReadFile("package.json")
	if err != nil {
		return nil, fmt.Errorf("failed to read package.json: %w", err)
	}

	pkg, err := ParsePackageJSON(pkgData)
	if err != nil {
		return nil, fmt.Errorf("failed to parse package.json: %w", err)
	}

	// Detect package manager
	pmInfo := DetectPackageManager(ctx, pkg)

	// Detect Node.js version
	nodeVersion := DetectNodeVersion(ctx, pkg)

	// Detect framework
	fwInfo := DetectFramework(ctx, pkg)

	// Build the plan
	language := "nodejs"
	languageVersion := nodeVersion
	if pmInfo.Name == PackageManagerBun {
		language = "bun"
		// For bun, use bun version as language version
		if pmInfo.Version != "" {
			languageVersion = pmInfo.Version
		} else {
			languageVersion = "latest"
		}
	}

	plan := &app.Plan{
		Provider:              "node",
		Language:              language,
		LanguageVersion:       languageVersion,
		PackageManager:        string(pmInfo.Name),
		PackageManagerVersion: pmInfo.Version,
		DetectedFiles:         []string{"package.json"},
		Metadata:              make(map[string]interface{}),
	}

	// Add runtime info for bun
	if pmInfo.Name == PackageManagerBun {
		plan.Metadata["runtime"] = "bun"
		plan.Metadata["runtime_note"] = "Using Bun runtime (oven/bun image)"
	}

	// Add framework info
	if fwInfo.Name != FrameworkNone {
		plan.Framework = string(fwInfo.Name)
		plan.FrameworkVersion = fwInfo.Version
		if fwInfo.OutputType != OutputTypeNone {
			plan.Metadata["output_type"] = string(fwInfo.OutputType)
		}
	}

	// Determine install command
	plan.InstallCommand = pmInfo.GetInstallCommand()

	// Determine build command
	plan.BuildCommand = determineBuildCommand(pkg, pmInfo, fwInfo)

	// Determine start command
	plan.StartCommand = determineStartCommand(pkg, pmInfo, fwInfo)

	// Add detected files to the list
	plan.DetectedFiles = append(plan.DetectedFiles, detectRelevantFiles(ctx, pmInfo)...)

	// Add additional metadata
	if pkg.Name != "" {
		plan.Metadata["name"] = pkg.Name
	}
	if pkg.Version != "" {
		plan.Metadata["version"] = pkg.Version
	}
	if pkg.IsMonorepo() {
		plan.Metadata["is_monorepo"] = true
		plan.Metadata["workspaces"] = pkg.Workspaces.Packages
	}
	if pkg.Type != "" {
		plan.Metadata["module_type"] = pkg.Type
	}

	// Detect native dependencies
	nativeDeps := DetectNativeDependencies(pkg)
	if len(nativeDeps) > 0 {
		aptPackages := GetRequiredAptPackages(nativeDeps)
		plan.Metadata["apt_packages"] = aptPackages

		// Track which native packages were detected
		var detected []string
		for _, dep := range nativeDeps {
			detected = append(detected, dep.Package)
		}
		plan.Metadata["native_packages"] = detected
	}

	// Check for base image override
	if baseImage := ctx.Env["COOLPACK_BASE_IMAGE"]; baseImage != "" {
		plan.Metadata["base_image"] = baseImage
	}

	// Detect Cypress for cache
	if pkg.HasDependency("cypress") {
		plan.Metadata["has_cypress"] = true
	}

	// Detect moon repo
	if ctx.HasFile(".moon/workspace.yml") {
		plan.Metadata["has_moon"] = true
	}

	// Custom cache directories from package.json
	if len(pkg.CacheDirectories) > 0 {
		plan.Metadata["cache_directories"] = pkg.CacheDirectories
	}

	// Detect SPA (only for static output)
	if outputType := plan.Metadata["output_type"]; outputType == "static" {
		if isSPA := detectSPA(pkg, fwInfo); isSPA {
			plan.Metadata["is_spa"] = true
		}
	}

	return plan, nil
}

// detectSPA checks if the application is a Single Page Application
// by looking for client-side router dependencies
func detectSPA(pkg *PackageJSON, fw FrameworkInfo) bool {
	// Frameworks that handle routing server-side or generate static HTML per route
	// don't need SPA fallback even in static mode
	switch fw.Name {
	case FrameworkGatsby, FrameworkEleventy:
		// Static site generators that create HTML for each route
		return false
	case FrameworkNextJS, FrameworkNuxt, FrameworkAstro:
		// These generate static HTML per route in export mode
		return false
	}

	// Client-side router dependencies indicate SPA
	spaRouters := []string{
		// Vue
		"vue-router",
		// React
		"react-router-dom",
		"react-router",
		"@reach/router",
		"wouter",
		"@tanstack/react-router",
		// Svelte (not SvelteKit)
		"svelte-navigator",
		"svelte-routing",
		"@roxi/routify",
		// Solid
		"@solidjs/router",
		"solid-app-router",
		// Preact
		"preact-router",
		// General
		"navigo",
		"page",
	}

	for _, router := range spaRouters {
		if pkg.HasDependency(router) {
			return true
		}
	}

	return false
}

// determineBuildCommand determines the build command to use
func determineBuildCommand(pkg *PackageJSON, pm PackageManagerInfo, fw FrameworkInfo) string {
	run := pm.GetRunCommand()

	// Check for explicit build script
	if pkg.HasScript("build") {
		return run + " build"
	}

	// Use framework-specific defaults
	if cmd := fw.GetDefaultBuildCommand(pm); cmd != "" {
		return cmd
	}

	return ""
}

// determineStartCommand determines the start command to use
func determineStartCommand(pkg *PackageJSON, pm PackageManagerInfo, fw FrameworkInfo) string {
	run := pm.GetRunCommand()

	// Check for explicit start script
	if pkg.HasScript("start") {
		return run + " start"
	}

	// Check for explicit serve script (common for SPAs)
	if pkg.HasScript("serve") {
		return run + " serve"
	}

	// Use framework-specific defaults
	if cmd := fw.GetDefaultStartCommand(pm); cmd != "" {
		return cmd
	}

	// Fallback: check main entry point
	if pkg.Main != "" {
		return fmt.Sprintf("node %s", pkg.Main)
	}

	// Check for common entry points
	entryPoints := []string{"dist/index.js", "build/index.js", "index.js", "server.js", "app.js"}
	for _, ep := range entryPoints {
		if hasEntryPoint(pkg, ep) {
			return fmt.Sprintf("node %s", ep)
		}
	}

	return ""
}

// hasEntryPoint checks if the entry point might exist (based on package.json hints)
func hasEntryPoint(pkg *PackageJSON, path string) bool {
	// This is a simple heuristic - in a real implementation we might check the filesystem
	if pkg.Main == path {
		return true
	}
	return false
}

// detectRelevantFiles returns a list of relevant files that were detected
func detectRelevantFiles(ctx *app.Context, pm PackageManagerInfo) []string {
	var files []string

	// Lock file
	lockFile := pm.GetLockFile()
	if ctx.HasFile(lockFile) {
		files = append(files, lockFile)
	}

	// Version files
	versionFiles := []string{".nvmrc", ".node-version", ".tool-versions", "mise.toml"}
	for _, f := range versionFiles {
		if ctx.HasFile(f) {
			files = append(files, f)
		}
	}

	// Config files
	configFiles := []string{
		".yarnrc.yml", ".yarnrc.yaml", ".npmrc", ".pnpmrc",
		"tsconfig.json", "jsconfig.json",
		"vite.config.js", "vite.config.ts", "vite.config.mjs",
		"next.config.js", "next.config.mjs", "next.config.ts",
		"astro.config.mjs", "astro.config.js", "astro.config.ts",
		"angular.json",
		"remix.config.js",
		"nuxt.config.ts", "nuxt.config.js",
	}
	for _, f := range configFiles {
		if ctx.HasFile(f) {
			files = append(files, f)
		}
	}

	return files
}
