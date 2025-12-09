package app

// Plan represents the detected build plan for an application
type Plan struct {
	// Provider is the name of the provider that detected this application (e.g., "node")
	Provider string `json:"provider"`

	// Language is the detected programming language (e.g., "nodejs", "python")
	Language string `json:"language"`

	// LanguageVersion is the detected or default version of the language runtime
	LanguageVersion string `json:"language_version,omitempty"`

	// Framework is the detected framework (e.g., "nextjs", "remix", "express")
	Framework string `json:"framework,omitempty"`

	// FrameworkVersion is the version of the detected framework
	FrameworkVersion string `json:"framework_version,omitempty"`

	// PackageManager is the detected package manager (e.g., "npm", "yarn", "pnpm", "bun")
	PackageManager string `json:"package_manager,omitempty"`

	// PackageManagerVersion is the version of the package manager
	PackageManagerVersion string `json:"package_manager_version,omitempty"`

	// InstallCommand is the command to install dependencies
	InstallCommand string `json:"install_command,omitempty"`

	// BuildCommand is the command to build the application
	BuildCommand string `json:"build_command,omitempty"`

	// StartCommand is the command to start the application
	StartCommand string `json:"start_command,omitempty"`

	// DetectedFiles lists the files that were used for detection
	DetectedFiles []string `json:"detected_files,omitempty"`

	// Metadata contains additional provider-specific information
	Metadata map[string]interface{} `json:"metadata,omitempty"`

	// BuildEnv contains environment variables available during build (ARG in Dockerfile)
	BuildEnv map[string]string `json:"build_env,omitempty"`

	// Env contains environment variables available at runtime (ENV in Dockerfile)
	Env map[string]string `json:"env,omitempty"`
}
