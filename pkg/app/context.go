package app

import (
	"os"
	"path/filepath"
)

// Context provides information about the application being analyzed
type Context struct {
	// Path is the absolute path to the application root
	Path string

	// Env contains environment variables that may influence detection
	Env map[string]string
}

// NewContext creates a new Context for the given path
func NewContext(path string) *Context {
	return &Context{
		Path: path,
		Env:  make(map[string]string),
	}
}

// HasFile checks if a file exists in the application path
func (ctx *Context) HasFile(name string) bool {
	path := filepath.Join(ctx.Path, name)
	_, err := os.Stat(path)
	return err == nil
}

// ReadFile reads a file from the application path
func (ctx *Context) ReadFile(name string) ([]byte, error) {
	path := filepath.Join(ctx.Path, name)
	return os.ReadFile(path)
}

// ListFiles lists files matching a pattern in the application path
func (ctx *Context) ListFiles(pattern string) ([]string, error) {
	fullPattern := filepath.Join(ctx.Path, pattern)
	matches, err := filepath.Glob(fullPattern)
	if err != nil {
		return nil, err
	}

	// Convert to relative paths
	result := make([]string, 0, len(matches))
	for _, m := range matches {
		rel, err := filepath.Rel(ctx.Path, m)
		if err != nil {
			continue
		}
		result = append(result, rel)
	}

	return result, nil
}
