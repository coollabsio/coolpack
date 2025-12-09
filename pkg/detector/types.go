package detector

import (
	"github.com/coollabsio/coolpack/pkg/app"
)

// Plan is an alias to app.Plan for convenience
type Plan = app.Plan

// Provider is the interface that all language/framework providers must implement
type Provider interface {
	// Name returns the name of the provider
	Name() string

	// Detect checks if this provider can handle the application at the given path
	// Returns true if the provider detected a matching application
	Detect(ctx *app.Context) (bool, error)

	// Plan generates a build plan for the detected application
	Plan(ctx *app.Context) (*app.Plan, error)
}
