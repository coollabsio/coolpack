# Coolpack

A general-purpose build pack that automatically detects your application type, generates optimized Dockerfiles, and builds production-ready container images.

## Features

- **Auto-detection** - Automatically detects language, framework, and package manager
- **Optimized Dockerfiles** - Generates multi-stage builds with security best practices
- **BuildKit caching** - Framework-specific cache mounts for faster rebuilds
- **Non-root containers** - Runs as unprivileged user by default
- **Native dependency support** - Automatically installs required system packages
- **Static site support** - Serves static builds with Caddy (default) or nginx

## Supported Frameworks

More coming soon, we will support a bunch of them, not just Javascript based.

| Framework | Output Types |
|-----------|--------------|
| Next.js | Server (SSR) / Static |
| Nuxt | Server (SSR) / Static |
| Remix / React Router | Server (SSR) / Static |
| Astro | Server (SSR) / Static |
| SvelteKit | Server (SSR) / Static |
| Solid Start | Server (SSR) / Static |
| TanStack Start | Server (SSR) / Static |
| Vite | Static |
| Gatsby | Static |
| Angular | Server (SSR) / Static |
| Express | Server |
| Fastify | Server |
| NestJS | Server |
| AdonisJS | Server |

## Installation

### Quick Install

```bash
curl -fsSL https://raw.githubusercontent.com/coollabsio/coolpack/main/install.sh | bash
```

### From Source

```bash
git clone https://github.com/coollabsio/coolpack.git
cd coolpack
./build.sh
```

The binary will be created at `./coolpack`.

### Requirements

- Docker with BuildKit support (for building images)
- Go 1.21+ (only for building from source)

## Quick Start

```bash
# Navigate to your project
cd my-app

# See what Coolpack detects
coolpack plan

# Generate Dockerfile
coolpack prepare

# Build container image
coolpack build

# Run the container (development only)
coolpack run
```

## Commands

### `coolpack plan [path]`

Analyze and display the build plan without generating any files.

```bash
coolpack plan                    # Current directory
coolpack plan ./my-app           # Specific path
coolpack plan --json             # Output as JSON
coolpack plan --out              # Save to coolpack.json
coolpack plan --out custom.json  # Save to custom file
coolpack plan --packages curl --packages wget  # Add custom packages
coolpack plan --build-env NEXT_PUBLIC_API_URL=https://api.example.com  # Add build env
```

**Flags:**
| Flag | Description |
|------|-------------|
| `--json` | Output as JSON |
| `-o, --out` | Write plan to file (default: `coolpack.json`) |
| `--packages` | Additional APT packages to install |
| `--build-env` | Build-time env vars (KEY=value or KEY) |

### `coolpack prepare [path]`

Generate a Dockerfile in the `.coolpack/` directory.

If a `coolpack.json` file exists in the project root, it will be used instead of running detection.

```bash
coolpack prepare
coolpack prepare --static-server nginx     # Use nginx instead of Caddy
coolpack prepare --build-cmd "npm run build:prod"
coolpack prepare --plan coolpack.json      # Use specific plan file
coolpack prepare --packages curl           # Add custom APT packages
```

**Flags:**
| Flag | Description |
|------|-------------|
| `-i, --install-cmd` | Override install command |
| `-b, --build-cmd` | Override build command |
| `-s, --start-cmd` | Override start command |
| `--static-server` | Static server: `caddy` (default), `nginx` |
| `--output-dir` | Override static output directory (e.g., `dist`, `build`) |
| `--spa` | Enable SPA mode (serves index.html for all routes) |
| `--no-spa` | Disable SPA mode (overrides auto-detection) |
| `--build-env` | Build-time env vars (KEY=value or KEY) |
| `--packages` | Additional APT packages to install |
| `--plan` | Use plan file instead of detection |

### `coolpack build [path]`

Generate Dockerfile and build the container image.

If a `coolpack.json` file exists in the project root, it will be used instead of running detection.

```bash
coolpack build
coolpack build -n my-app -t v1.0.0
coolpack build --build-env NEXT_PUBLIC_API_URL=https://api.example.com
coolpack build --no-cache
coolpack build --plan coolpack.json        # Use specific plan file
coolpack build --packages ffmpeg           # Add custom APT packages
```

**Flags:**
| Flag | Description |
|------|-------------|
| `-n, --name` | Image name (defaults to directory name) |
| `-t, --tag` | Image tag (default: `latest`) |
| `--no-cache` | Build without Docker cache |
| `-i, --install-cmd` | Override install command |
| `-b, --build-cmd` | Override build command |
| `-s, --start-cmd` | Override start command |
| `--static-server` | Static server: `caddy` (default), `nginx` |
| `--output-dir` | Override static output directory (e.g., `dist`, `build`) |
| `--spa` | Enable SPA mode (serves index.html for all routes) |
| `--no-spa` | Disable SPA mode (overrides auto-detection) |
| `--build-env` | Build-time env vars |
| `--packages` | Additional APT packages to install |
| `--plan` | Use plan file instead of detection |

### `coolpack run [path]`

Build and run the container locally. **For development only.**

```bash
coolpack run
coolpack run -e DATABASE_URL=postgres://localhost/db
```

**Flags:**
| Flag | Description |
|------|-------------|
| `-n, --name` | Image name |
| `-t, --tag` | Image tag |
| `-e, --env` | Runtime env vars (KEY=value) |

### `coolpack version`

Print version information.

```bash
coolpack version
```

## Configuration

### Environment Variables

Override Coolpack behavior with environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `COOLPACK_INSTALL_CMD` | Override install command | Auto-detected |
| `COOLPACK_BUILD_CMD` | Override build command | Auto-detected |
| `COOLPACK_START_CMD` | Override start command | Auto-detected |
| `COOLPACK_BASE_IMAGE` | Override base Docker image | Provider-specific |
| `COOLPACK_NODE_VERSION` | Override Node.js version | Auto-detected or `24` |
| `COOLPACK_STATIC_SERVER` | Static file server | `caddy` |
| `COOLPACK_SPA_OUTPUT_DIR` | Override static output directory | Framework-specific |
| `COOLPACK_SPA` | Enable SPA mode | Auto-detected |
| `COOLPACK_NO_SPA` | Disable SPA mode | `false` |
| `COOLPACK_PACKAGES` | Additional APT packages (comma-separated) | - |
| `NODE_VERSION` | Alternative to `COOLPACK_NODE_VERSION` (legacy) | - |

**Priority:** CLI flags > Environment variables > Auto-detected

**Default Base Images by Provider:**
| Provider | Default Base Image |
|----------|-------------------|
| Node.js | `node:<version>-slim` |
| Node.js (bun) | `oven/bun:<version>-slim` |

### Build-time vs Runtime Environment Variables

**Build-time** variables are baked into the image during build:

```bash
coolpack build --build-env NEXT_PUBLIC_API_URL=https://api.example.com
```

Use for:
- Next.js `NEXT_PUBLIC_*` variables
- Vite `VITE_*` variables
- Any `process.env` accessed during build

**Runtime** variables are passed when running the container:

```bash
docker run -e DATABASE_URL=postgres://... -e API_KEY=... my-app:latest
```

Use for:
- Secrets that shouldn't be in the image
- Config that changes per environment

### Custom Cache Directories

Add custom cache directories in `package.json`:

```json
{
  "cacheDirectories": [
    ".cache",
    "tmp/build-cache"
  ]
}
```

### Node.js Version

Coolpack detects Node.js version from (in priority order):

1. `COOLPACK_NODE_VERSION` env var
2. `engines.node` in package.json
3. `.nvmrc` file
4. `.node-version` file
5. `.tool-versions` file (asdf)
6. `mise.toml` file
7. Default: `24`

### Package Manager

Detected from (in priority order):

1. `packageManager` field in package.json
2. Lock files (`pnpm-lock.yaml`, `bun.lockb`, `yarn.lock`, `package-lock.json`)
3. `engines` field in package.json
4. Default: `npm`

## Examples

### Next.js with SSR

```bash
cd my-nextjs-app
coolpack build -n my-app -t latest
docker run -p 3000:3000 my-app:latest
```

### Static Vite App

```bash
cd my-vite-app
coolpack build
docker run -p 80:80 my-vite-app:latest
```

### Custom Build Commands

```bash
coolpack build \
  --install-cmd "pnpm install --frozen-lockfile" \
  --build-cmd "pnpm build:production" \
  --start-cmd "node dist/server.js"
```

### With Build-time Variables

```bash
coolpack build \
  --build-env NEXT_PUBLIC_API_URL=https://api.example.com \
  --build-env NEXT_PUBLIC_GA_ID=UA-123456
```

### Using nginx for Static Sites

```bash
coolpack build --static-server nginx
```

### SPA with Client-side Routing

For Single Page Applications, Coolpack auto-detects client-side routers (vue-router, react-router-dom, etc.) and configures the server to serve `index.html` for all routes:

```bash
# Auto-detected when vue-router, react-router-dom, etc. is found
coolpack build

# Manual override
coolpack build --spa
```

### Using Plan Files

Save a build plan and reuse it for reproducible builds:

```bash
# Generate and save plan
coolpack plan --out

# Edit coolpack.json to customize settings...

# Build using the plan file (auto-detected)
coolpack build

# Or specify explicitly
coolpack build --plan coolpack.json
```

### Custom APT Packages

Add system packages that aren't auto-detected:

```bash
# Via CLI flag
coolpack build --packages ffmpeg --packages curl

# Via environment variable
COOLPACK_PACKAGES=ffmpeg,curl coolpack build

# Via plan file (add to metadata.custom_packages)
```

---

## Development

### Prerequisites

- Go 1.21+
- Docker (for testing builds)

### Setup

```bash
git clone https://github.com/coollabsio/coolpack.git
cd coolpack
go mod download
```

### Building

```bash
./build.sh
```

### Project Structure

```
coolpack/
├── main.go                          # Entry point
├── build.sh                         # Build script
├── cmd/coolpack/
│   ├── root.go                      # Root CLI command
│   ├── plan.go                      # Plan subcommand
│   ├── prepare.go                   # Prepare subcommand
│   ├── build.go                     # Build subcommand
│   └── run.go                       # Run subcommand
└── pkg/
    ├── app/
    │   ├── context.go               # App context (path, env, file helpers)
    │   └── plan.go                  # Plan struct
    ├── detector/
    │   ├── detector.go              # Main detector, registers providers
    │   └── types.go                 # Provider interface
    ├── generator/
    │   └── generator.go             # Dockerfile generation
    └── providers/node/
        ├── node.go                  # Node.js provider
        ├── package_json.go          # package.json parsing
        ├── package_manager.go       # Package manager detection
        ├── version.go               # Node version detection
        ├── framework.go             # Framework detection
        ├── config_parser.go         # JS/TS config parsing
        └── native_deps.go           # Native dependency detection
```

### Adding a New Provider

1. Create `pkg/providers/<name>/<name>.go`
2. Implement the `Provider` interface:

```go
type Provider interface {
    Name() string
    Detect(ctx *app.Context) (bool, error)
    Plan(ctx *app.Context) (*app.Plan, error)
}
```

3. Register in `pkg/detector/detector.go`:

```go
func (d *Detector) registerProviders() {
    d.providers = append(d.providers, node.New())
    d.providers = append(d.providers, yourprovider.New())
}
```

### Adding Framework Detection

For Node.js frameworks, edit `pkg/providers/node/framework.go`:

1. Add framework constant
2. Add detection logic in `DetectFramework()`
3. Add output type determination
4. Add default build/start commands

### Testing

Test against example projects:

```bash
# Test detection
./coolpack plan /path/to/test-project

# Test Dockerfile generation
./coolpack prepare /path/to/test-project
cat /path/to/test-project/.coolpack/Dockerfile

# Test full build
./coolpack build /path/to/test-project
```

### Dependencies

- `github.com/spf13/cobra` - CLI framework
- `github.com/smacker/go-tree-sitter` - AST parsing for JS/TS config files

---

## How It Works

1. **Detection** - Scans project files to identify language, framework, and package manager
2. **Planning** - Creates a build plan with install, build, and start commands
3. **Generation** - Produces an optimized multi-stage Dockerfile
4. **Building** - Runs `docker build` with BuildKit cache mounts

### Generated Dockerfile Features

- Multi-stage builds (builder + runner)
- BuildKit cache mounts for dependencies and build artifacts
- Non-root user (`cooluser`, UID 1001)
- Production-optimized Node.js settings
- Framework-specific output copying
- Automatic native dependency installation

## License

MIT

## Contributing

Contributions are welcome! Please feel free to submit issues and pull requests.
