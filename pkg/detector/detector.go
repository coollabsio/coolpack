package detector

import (
	"os"

	"github.com/coollabsio/coolpack/pkg/app"
	"github.com/coollabsio/coolpack/pkg/providers/node"
)

// Detector handles application detection using registered providers
type Detector struct {
	path      string
	providers []Provider
}

// New creates a new Detector for the given path
func New(path string) *Detector {
	d := &Detector{
		path:      path,
		providers: make([]Provider, 0),
	}

	// Register all providers
	d.registerProviders()

	return d
}

// registerProviders adds all available providers to the detector
func (d *Detector) registerProviders() {
	// Node.js provider
	d.providers = append(d.providers, node.New())

	// TODO: Add more providers here (python, go, rust, etc.)
}

// Detect runs detection using all registered providers and returns a plan
func (d *Detector) Detect() (*Plan, error) {
	ctx := app.NewContext(d.path)

	// Load environment variables that might influence detection
	ctx.Env = loadRelevantEnvVars()

	// Try each provider in order
	for _, provider := range d.providers {
		detected, err := provider.Detect(ctx)
		if err != nil {
			// Log error but continue to next provider
			continue
		}

		if detected {
			return provider.Plan(ctx)
		}
	}

	return nil, nil
}

// loadRelevantEnvVars loads environment variables that influence detection
func loadRelevantEnvVars() map[string]string {
	env := make(map[string]string)

	// Coolpack config
	envVars := []string{
		// Command overrides
		"COOLPACK_INSTALL_CMD",
		"COOLPACK_BUILD_CMD",
		"COOLPACK_START_CMD",
		// Image and version overrides
		"COOLPACK_BASE_IMAGE",
		"COOLPACK_NODE_VERSION",
		"COOLPACK_SPA_OUTPUT_DIR",
		// Static server (caddy or nginx)
		"COOLPACK_STATIC_SERVER",
		// Legacy support
		"NODE_VERSION",
	}

	for _, v := range envVars {
		if val := os.Getenv(v); val != "" {
			env[v] = val
		}
	}

	return env
}
