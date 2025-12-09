package node

import (
	"github.com/coollabsio/coolpack/pkg/app"
	sitter "github.com/smacker/go-tree-sitter"
)

// Framework represents a detected Node.js framework
type Framework string

const (
	FrameworkNone        Framework = ""
	FrameworkNextJS      Framework = "nextjs"
	FrameworkRemix       Framework = "remix"
	FrameworkNuxt        Framework = "nuxt"
	FrameworkAstro       Framework = "astro"
	FrameworkVite        Framework = "vite"
	FrameworkCRA         Framework = "create-react-app"
	FrameworkAngular     Framework = "angular"
	FrameworkSvelteKit   Framework = "sveltekit"
	FrameworkSolidStart  Framework = "solid-start"
	FrameworkExpress     Framework = "express"
	FrameworkFastify     Framework = "fastify"
	FrameworkNestJS      Framework = "nestjs"
	FrameworkAdonisJS    Framework = "adonisjs"
	FrameworkReactRouter Framework = "react-router"
	FrameworkTanStack    Framework = "tanstack-start"
	FrameworkGatsby      Framework = "gatsby"
	FrameworkEleventy    Framework = "eleventy"
)

// OutputType represents the type of output the framework produces
type OutputType string

const (
	OutputTypeNone   OutputType = ""       // Unknown
	OutputTypeStatic OutputType = "static" // Static files (SPA, SSG) - can be served from any static file server
	OutputTypeServer OutputType = "server" // Needs Node.js server at runtime (SSR frameworks, backend APIs)
)

// FrameworkInfo contains information about the detected framework
type FrameworkInfo struct {
	Name       Framework
	Version    string
	OutputType OutputType
}

// DetectFramework detects the framework used by the project
func DetectFramework(ctx *app.Context, pkg *PackageJSON) FrameworkInfo {
	info := FrameworkInfo{
		Name:       FrameworkNone,
		Version:    "",
		OutputType: OutputTypeNone,
	}

	if pkg == nil {
		return info
	}

	// Meta-frameworks with SSR (check these first as they're more specific)
	if pkg.HasDependency("next") {
		info.Name = FrameworkNextJS
		info.Version = cleanVersion(pkg.GetDependencyVersion("next"))
		// Check if it's a static export (output: 'export' in next.config.*)
		if isNextJSStaticExport(ctx) {
			info.OutputType = OutputTypeStatic
		} else {
			info.OutputType = OutputTypeServer
		}
		return info
	}

	if pkg.HasDependency("@remix-run/react") || pkg.HasDependency("@remix-run/node") {
		info.Name = FrameworkRemix
		info.Version = cleanVersion(pkg.GetDependencyVersion("@remix-run/react"))
		info.OutputType = OutputTypeServer
		return info
	}

	if pkg.HasDependency("nuxt") || pkg.HasDependency("nuxt3") {
		info.Name = FrameworkNuxt
		info.Version = cleanVersion(pkg.GetDependencyVersion("nuxt"))
		// Check for ssr: false in nuxt.config.*
		if isNuxtSPAMode(ctx) {
			info.OutputType = OutputTypeStatic
		} else {
			info.OutputType = OutputTypeServer
		}
		return info
	}

	if pkg.HasDependency("astro") || ctx.HasFile("astro.config.mjs") || ctx.HasFile("astro.config.js") || ctx.HasFile("astro.config.ts") {
		info.Name = FrameworkAstro
		info.Version = cleanVersion(pkg.GetDependencyVersion("astro"))
		// Astro is static by default, SSR requires output: 'server' or 'hybrid'
		if isAstroSSRMode(ctx) {
			info.OutputType = OutputTypeServer
		} else {
			info.OutputType = OutputTypeStatic
		}
		return info
	}

	if pkg.HasDependency("@sveltejs/kit") {
		info.Name = FrameworkSvelteKit
		info.Version = cleanVersion(pkg.GetDependencyVersion("@sveltejs/kit"))
		// Check if using static adapter
		if pkg.HasDependency("@sveltejs/adapter-static") {
			info.OutputType = OutputTypeStatic
		} else {
			info.OutputType = OutputTypeServer
		}
		return info
	}

	if pkg.HasDependency("solid-start") || pkg.HasDependency("@solidjs/start") {
		info.Name = FrameworkSolidStart
		// Try @solidjs/start first (newer), then solid-start (older)
		version := pkg.GetDependencyVersion("@solidjs/start")
		if version == "" {
			version = pkg.GetDependencyVersion("solid-start")
		}
		info.Version = cleanVersion(version)
		// Check for ssr: false in app.config.*
		if isSolidStartSPAMode(ctx) {
			info.OutputType = OutputTypeStatic
		} else {
			info.OutputType = OutputTypeServer
		}
		return info
	}

	if pkg.HasDependency("@tanstack/start") || pkg.HasDependency("@tanstack/react-start") {
		info.Name = FrameworkTanStack
		info.Version = cleanVersion(pkg.GetDependencyVersion("@tanstack/start"))
		// Check for server.preset: 'static' in app.config.*
		if isTanStackStartStaticMode(ctx) {
			info.OutputType = OutputTypeStatic
		} else {
			info.OutputType = OutputTypeServer
		}
		return info
	}

	// React Router v7+ with config file is Remix
	if pkg.HasDependency("react-router") && (ctx.HasFile("react-router.config.ts") || ctx.HasFile("react-router.config.js")) {
		info.Name = FrameworkRemix
		info.Version = cleanVersion(pkg.GetDependencyVersion("react-router"))
		// Check for ssr: false in react-router.config.*
		if isReactRouterSPAMode(ctx) {
			info.OutputType = OutputTypeStatic
		} else {
			info.OutputType = OutputTypeServer
		}
		return info
	}

	if pkg.HasDependency("gatsby") {
		info.Name = FrameworkGatsby
		info.Version = cleanVersion(pkg.GetDependencyVersion("gatsby"))
		info.OutputType = OutputTypeStatic
		return info
	}

	if pkg.HasDependency("@11ty/eleventy") {
		info.Name = FrameworkEleventy
		info.Version = cleanVersion(pkg.GetDependencyVersion("@11ty/eleventy"))
		info.OutputType = OutputTypeStatic
		return info
	}

	// Angular detection (check before backend frameworks since Angular SSR uses Express)
	if pkg.HasDependency("@angular/core") || ctx.HasFile("angular.json") {
		info.Name = FrameworkAngular
		info.Version = cleanVersion(pkg.GetDependencyVersion("@angular/core"))
		// Check for @angular/ssr for SSR mode
		if pkg.HasDependency("@angular/ssr") {
			info.OutputType = OutputTypeServer
		} else {
			info.OutputType = OutputTypeStatic
		}
		return info
	}

	// Backend frameworks (need Node.js server at runtime)
	if pkg.HasDependency("@adonisjs/core") {
		info.Name = FrameworkAdonisJS
		info.Version = cleanVersion(pkg.GetDependencyVersion("@adonisjs/core"))
		info.OutputType = OutputTypeServer
		return info
	}

	if pkg.HasDependency("@nestjs/core") {
		info.Name = FrameworkNestJS
		info.Version = cleanVersion(pkg.GetDependencyVersion("@nestjs/core"))
		info.OutputType = OutputTypeServer
		return info
	}

	if pkg.HasDependency("fastify") {
		info.Name = FrameworkFastify
		info.Version = cleanVersion(pkg.GetDependencyVersion("fastify"))
		info.OutputType = OutputTypeServer
		return info
	}

	if pkg.HasDependency("express") {
		info.Name = FrameworkExpress
		info.Version = cleanVersion(pkg.GetDependencyVersion("express"))
		info.OutputType = OutputTypeServer
		return info
	}

	// Check for Create React App
	if pkg.HasDependency("react-scripts") {
		info.Name = FrameworkCRA
		info.Version = cleanVersion(pkg.GetDependencyVersion("react-scripts"))
		info.OutputType = OutputTypeStatic
		return info
	}

	// Vite detection (check after more specific frameworks)
	// Vite always produces static output
	if pkg.HasDependency("vite") || ctx.HasFile("vite.config.js") || ctx.HasFile("vite.config.ts") || ctx.HasFile("vite.config.mjs") {
		info.Name = FrameworkVite
		info.Version = cleanVersion(pkg.GetDependencyVersion("vite"))
		info.OutputType = OutputTypeStatic
		return info
	}

	return info
}

// isSPAProject checks if the project appears to be a Single Page Application
func isSPAProject(pkg *PackageJSON) bool {
	spaIndicators := []string{
		"react",
		"vue",
		"svelte",
		"preact",
		"lit",
		"solid-js",
		"@builder.io/qwik",
	}

	for _, dep := range spaIndicators {
		if pkg.HasDependency(dep) {
			return true
		}
	}
	return false
}

// isNextJSStaticExport checks if Next.js is configured for static export
// by looking for output: 'export' in next.config.* files using tree-sitter
func isNextJSStaticExport(ctx *app.Context) bool {
	parser := NewConfigParser()

	// Check TypeScript config first
	if ctx.HasFile("next.config.ts") {
		data, err := ctx.ReadFile("next.config.ts")
		if err == nil {
			root, err := parser.ParseTS(data)
			if err == nil {
				value := FindPropertyValue(root, data, "output")
				if value == "export" {
					return true
				}
			}
		}
	}

	// Check JavaScript configs
	jsConfigs := []string{"next.config.mjs", "next.config.js"}
	for _, configFile := range jsConfigs {
		if ctx.HasFile(configFile) {
			data, err := ctx.ReadFile(configFile)
			if err != nil {
				continue
			}
			root, err := parser.ParseJS(data)
			if err != nil {
				continue
			}
			value := FindPropertyValue(root, data, "output")
			if value == "export" {
				return true
			}
		}
	}

	return false
}

// isAstroSSRMode checks if Astro is configured for SSR mode
// by looking for output: 'server' or 'hybrid' in astro.config.* files
func isAstroSSRMode(ctx *app.Context) bool {
	parser := NewConfigParser()

	configFiles := []string{"astro.config.ts", "astro.config.mjs", "astro.config.js"}
	for _, configFile := range configFiles {
		if ctx.HasFile(configFile) {
			data, err := ctx.ReadFile(configFile)
			if err != nil {
				continue
			}

			var root *sitter.Node
			if configFile == "astro.config.ts" {
				root, err = parser.ParseTS(data)
			} else {
				root, err = parser.ParseJS(data)
			}
			if err != nil {
				continue
			}

			value := FindPropertyValue(root, data, "output")
			if value == "server" || value == "hybrid" {
				return true
			}
		}
	}

	return false
}

// isNuxtSPAMode checks if Nuxt is configured for SPA mode
// by looking for ssr: false in nuxt.config.* files
func isNuxtSPAMode(ctx *app.Context) bool {
	parser := NewConfigParser()

	configFiles := []string{"nuxt.config.ts", "nuxt.config.js", "nuxt.config.mjs"}
	for _, configFile := range configFiles {
		if ctx.HasFile(configFile) {
			data, err := ctx.ReadFile(configFile)
			if err != nil {
				continue
			}

			var root *sitter.Node
			if configFile == "nuxt.config.ts" {
				root, err = parser.ParseTS(data)
			} else {
				root, err = parser.ParseJS(data)
			}
			if err != nil {
				continue
			}

			value := FindPropertyValue(root, data, "ssr")
			if value == "false" {
				return true
			}
		}
	}

	return false
}

// isReactRouterSPAMode checks if React Router/Remix is configured for SPA mode
// by looking for ssr: false in react-router.config.* files
func isReactRouterSPAMode(ctx *app.Context) bool {
	parser := NewConfigParser()

	// Check TypeScript config first
	if ctx.HasFile("react-router.config.ts") {
		data, err := ctx.ReadFile("react-router.config.ts")
		if err == nil {
			root, err := parser.ParseTS(data)
			if err == nil {
				value := FindPropertyValue(root, data, "ssr")
				if value == "false" {
					return true
				}
			}
		}
	}

	// Check JavaScript config
	if ctx.HasFile("react-router.config.js") {
		data, err := ctx.ReadFile("react-router.config.js")
		if err == nil {
			root, err := parser.ParseJS(data)
			if err == nil {
				value := FindPropertyValue(root, data, "ssr")
				if value == "false" {
					return true
				}
			}
		}
	}

	return false
}

// isSolidStartSPAMode checks if Solid Start is configured for SPA mode
// by looking for ssr: false in app.config.* files
func isSolidStartSPAMode(ctx *app.Context) bool {
	parser := NewConfigParser()

	// Check TypeScript config first
	if ctx.HasFile("app.config.ts") {
		data, err := ctx.ReadFile("app.config.ts")
		if err == nil {
			root, err := parser.ParseTS(data)
			if err == nil {
				value := FindPropertyValue(root, data, "ssr")
				if value == "false" {
					return true
				}
			}
		}
	}

	// Check JavaScript config
	if ctx.HasFile("app.config.js") {
		data, err := ctx.ReadFile("app.config.js")
		if err == nil {
			root, err := parser.ParseJS(data)
			if err == nil {
				value := FindPropertyValue(root, data, "ssr")
				if value == "false" {
					return true
				}
			}
		}
	}

	return false
}

// isTanStackStartStaticMode checks if TanStack Start is configured for static mode
// by looking for server.preset: 'static' in app.config.* files
func isTanStackStartStaticMode(ctx *app.Context) bool {
	parser := NewConfigParser()

	// Check TypeScript config first
	if ctx.HasFile("app.config.ts") {
		data, err := ctx.ReadFile("app.config.ts")
		if err == nil {
			root, err := parser.ParseTS(data)
			if err == nil {
				value := FindNestedPropertyValue(root, data, "server", "preset")
				if value == "static" {
					return true
				}
			}
		}
	}

	// Check JavaScript config
	if ctx.HasFile("app.config.js") {
		data, err := ctx.ReadFile("app.config.js")
		if err == nil {
			root, err := parser.ParseJS(data)
			if err == nil {
				value := FindNestedPropertyValue(root, data, "server", "preset")
				if value == "static" {
					return true
				}
			}
		}
	}

	return false
}

// cleanVersion removes common prefixes from version strings
func cleanVersion(v string) string {
	if len(v) > 0 && (v[0] == '^' || v[0] == '~' || v[0] == '>' || v[0] == '<' || v[0] == '=') {
		// Remove leading constraint characters
		for len(v) > 0 && (v[0] == '^' || v[0] == '~' || v[0] == '>' || v[0] == '<' || v[0] == '=' || v[0] == ' ') {
			v = v[1:]
		}
	}
	return v
}

// GetDefaultBuildCommand returns the default build command for a framework
func (f FrameworkInfo) GetDefaultBuildCommand(pm PackageManagerInfo) string {
	run := pm.GetRunCommand()

	switch f.Name {
	case FrameworkNextJS:
		return run + " build"
	case FrameworkRemix:
		return run + " build"
	case FrameworkNuxt:
		return run + " build"
	case FrameworkAstro:
		return run + " build"
	case FrameworkVite, FrameworkCRA:
		return run + " build"
	case FrameworkAngular:
		return run + " build"
	case FrameworkSvelteKit:
		return run + " build"
	case FrameworkGatsby:
		return run + " build"
	default:
		return ""
	}
}

// GetDefaultStartCommand returns the default start command for a framework
func (f FrameworkInfo) GetDefaultStartCommand(pm PackageManagerInfo) string {
	run := pm.GetRunCommand()

	switch f.Name {
	case FrameworkNextJS:
		return run + " start"
	case FrameworkRemix:
		return run + " start"
	case FrameworkNuxt:
		return "node .output/server/index.mjs"
	case FrameworkAstro:
		return "node ./dist/server/entry.mjs"
	case FrameworkNestJS, FrameworkExpress, FrameworkFastify:
		return run + " start"
	default:
		return ""
	}
}
