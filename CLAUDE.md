# Coolpack

> **Note**: Keep this file and `README.md` updated when making changes to detection logic, adding providers, adding CLI flags, or modifying the project structure.

A general purpose build pack for applications. Detects application type, generates Dockerfiles, and builds container images.

## Build

```bash
./build.sh
```

## Commands

- `coolpack plan [path]` - Detect and output build plan
- `coolpack plan [path] --json` - Output as JSON
- `coolpack prepare [path]` - Generate Dockerfile in `.coolpack/` directory
  - `-i, --install-cmd` - Override install command
  - `-b, --build-cmd` - Override build command
  - `-s, --start-cmd` - Override start command
  - `--static-server` - Static file server: `caddy` (default), `nginx`
  - `--output-dir` - Override static output directory (e.g., `dist`, `build`, `out`)
  - `--spa` - Enable SPA mode (serves index.html for all routes)
  - `--no-spa` - Disable SPA mode (overrides auto-detection)
  - `--build-env` - Build-time environment variables (KEY=value or KEY to pull from current env)
- `coolpack build [path]` - Build container image
  - `-n, --name` - Image name (defaults to directory name)
  - `-t, --tag` - Image tag (default "latest")
  - `--no-cache` - Build without Docker cache
  - `-i, --install-cmd` - Override install command
  - `-b, --build-cmd` - Override build command
  - `-s, --start-cmd` - Override start command
  - `--static-server` - Static file server: `caddy` (default), `nginx`
  - `--output-dir` - Override static output directory (e.g., `dist`, `build`, `out`)
  - `--spa` - Enable SPA mode (serves index.html for all routes)
  - `--no-spa` - Disable SPA mode (overrides auto-detection)
  - `--build-env` - Build-time environment variables (KEY=value or KEY to pull from current env)
- `coolpack run [path]` - Run container (**DEVELOPMENT ONLY**)
  - `-n, --name` - Image name (defaults to directory name)
  - `-t, --tag` - Image tag (default "latest")
  - `-e, --env` - Runtime environment variables (KEY=value)
- `coolpack version` - Print version information

## Environment Variables

Coolpack behavior can be configured via environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `COOLPACK_INSTALL_CMD` | Override install command | Auto-detected |
| `COOLPACK_BUILD_CMD` | Override build command | Auto-detected |
| `COOLPACK_START_CMD` | Override start command | Auto-detected |
| `COOLPACK_BASE_IMAGE` | Override the base Docker image (e.g., `node:20-alpine`) | Provider-specific |
| `COOLPACK_NODE_VERSION` | Override Node.js version | Auto-detected or `24` |
| `COOLPACK_STATIC_SERVER` | Static file server for static sites | `caddy` |
| `COOLPACK_SPA` | Enable SPA mode (serves index.html for all routes) | Auto-detected |
| `COOLPACK_NO_SPA` | Disable SPA mode (overrides auto-detection) | `false` |
| `COOLPACK_SPA_OUTPUT_DIR` | Override static output directory | Framework-specific |
| `NODE_VERSION` | Alternative to `COOLPACK_NODE_VERSION` (legacy) | - |

**Priority**: CLI flags > Environment variables > Auto-detected

**Default Base Images by Provider**:
| Provider | Default Base Image |
|----------|-------------------|
| Node.js | `node:<version>-slim` |
| Node.js (bun) | `oven/bun:<version>-slim` |

Example:
```bash
# Override commands via CLI
coolpack build --install-cmd "npm install" --build-cmd "npm run build:prod" --start-cmd "node server.js"

# Override commands via env vars
COOLPACK_BUILD_CMD="npm run build:custom" coolpack build

# Use nginx instead of Caddy for static sites
coolpack build --static-server nginx
COOLPACK_STATIC_SERVER=nginx coolpack build

# Use full Node image (not slim) for native dependencies
COOLPACK_BASE_IMAGE=node:20 coolpack build

# Force specific Node version
COOLPACK_NODE_VERSION=18 coolpack build
```

## Detection

### Node.js Provider

**Detection**: Project has `package.json` in root.

#### Node Version Detection (priority order)

1. `COOLPACK_NODE_VERSION` environment variable
2. `NODE_VERSION` environment variable (legacy)
3. `engines.node` field in package.json
4. `.nvmrc` file
5. `.node-version` file
6. `.tool-versions` file (asdf format)
7. `mise.toml` file
8. Default: `24`

#### Package Manager Detection (priority order)

1. `packageManager` field in package.json (e.g., `"pnpm@8.0.0"`)
2. Lock files:
   - `pnpm-lock.yaml` → pnpm
   - `bun.lockb` or `bun.lock` → bun
   - `.yarnrc.yml` or `.yarnrc.yaml` → yarn berry (v2+)
   - `yarn.lock` → yarn v1
   - `package-lock.json` → npm
3. `engines` field in package.json (pnpm, bun, yarn)
4. Default: npm

#### Framework Detection

Detected via dependencies in package.json or config files:

| Framework | Detection | Output Type |
|-----------|-----------|-------------|
| Next.js | `next` dependency | `server` (default) or `static` if `output: 'export'` |
| Remix | `@remix-run/*` or `react-router` + config file | `server` (default) or `static` if `ssr: false` |
| Nuxt | `nuxt` or `nuxt3` dependency | `server` (default) or `static` if `ssr: false` |
| Astro | `astro` dependency | `static` (default) or `server` if `output: 'server'/'hybrid'` |
| SvelteKit | `@sveltejs/kit` dependency | `server` (default) or `static` if `@sveltejs/adapter-static` |
| Solid Start | `solid-start` or `@solidjs/start` | `server` (default) or `static` if `ssr: false` |
| TanStack Start | `@tanstack/start` or `@tanstack/react-start` | `server` (default) or `static` if `server.preset: 'static'` |
| Gatsby | `gatsby` dependency | `static` |
| Eleventy | `@11ty/eleventy` dependency | `static` |
| Create React App | `react-scripts` dependency | `static` |
| Angular | `@angular/core` or `angular.json` | `server` if `@angular/ssr`, otherwise `static` |
| Vite | `vite` dependency or config files | `static` |
| NestJS | `@nestjs/core` dependency | `server` |
| Fastify | `fastify` dependency | `server` |
| Express | `express` dependency | `server` |
| AdonisJS | `@adonisjs/core` dependency | `server` |

#### Output Types

The `output_type` metadata field indicates how the application should be deployed:

| Type | Description |
|------|-------------|
| `static` | Static files - can be served from any static file server (Nginx, S3, CDN) |
| `server` | Needs Node.js server at runtime (SSR frameworks, backend APIs) |

#### Install Commands

| Package Manager | Command |
|-----------------|---------|
| npm | `npm ci` |
| yarn v1 | `yarn install --frozen-lockfile` |
| yarn berry | `yarn install --immutable` |
| pnpm | `pnpm install --frozen-lockfile` |
| bun | `bun install --frozen-lockfile` |

#### Build/Start Commands

- Uses `scripts.build` and `scripts.start` from package.json if present
- Falls back to framework-specific defaults
- Falls back to `main` field in package.json

#### Native Dependencies

Coolpack detects npm packages that require native system libraries and automatically installs the required APT packages in the Dockerfile.

| Package | APT Packages | Description |
|---------|--------------|-------------|
| `sharp` | `libvips-dev` | Image processing |
| `@prisma/client`, `prisma` | `openssl` | Database ORM |
| `puppeteer` | `chromium`, `libnss3`, `libatk*`, etc. | Headless Chrome |
| `playwright` | Browser dependencies | Browser automation |
| `canvas` | `libcairo2-dev`, `libjpeg-dev`, `libpango1.0-dev`, etc. | Canvas rendering |
| `bcrypt`, `argon2` | `build-essential`, `python3` | Password hashing |
| `sqlite3`, `better-sqlite3` | `build-essential`, `python3` | SQLite database |
| `node-gyp` | `build-essential`, `python3` | Native addon build tool |
| `ssh2` | `build-essential` | SSH client |
| `libsql`, `@libsql/client` | `build-essential` | LibSQL database |

If native dependencies aren't working with the slim image, override with full image:
```bash
COOLPACK_BASE_IMAGE=node:20 coolpack build
```

## Dockerfile Generation

The `prepare` command generates Dockerfiles in `.coolpack/` directory:

### Server Output (`output_type: "server"`)
- Multi-stage build with Node.js slim image
- Runs as non-root user `cooluser` (UID 1001)
- Exposes port 3000

### Static Output (`output_type: "static"`)
- Build stage with Node.js, serve stage with Caddy (default) or nginx
- Runs as non-root user `cooluser` (UID 1001)
- Exposes port 80
- Use `--static-server nginx` or `COOLPACK_STATIC_SERVER=nginx` to use nginx instead
- Framework-specific output directories (dist, out, build, etc.)

### SPA Mode

Single Page Applications need the server to serve `index.html` for all routes so the client-side router can handle them.

**Auto-detection**: Coolpack detects SPA mode when:
- Output type is `static`
- Project has a client-side router dependency:
  - `vue-router`
  - `react-router-dom`, `react-router`, `@reach/router`, `wouter`, `@tanstack/react-router`
  - `svelte-navigator`, `svelte-routing`, `@roxi/routify`
  - `@solidjs/router`, `solid-app-router`
  - `preact-router`

**Not auto-detected** for static site generators (Gatsby, Eleventy, Next.js export, Nuxt generate, Astro) that generate HTML for each route.

**Manual override**:
```bash
coolpack build --spa
COOLPACK_SPA=true coolpack build
```

When SPA mode is enabled:
- Caddy uses a Caddyfile with `try_files {path} /index.html`
- nginx uses `try_files $uri $uri/ /index.html`

### Build Environment Variables

Pass build-time environment variables with `--build-env` flag:
```bash
coolpack build --build-env API_URL=https://api.example.com --build-env DATABASE_URL
```

- `KEY=value` - Set explicit value
- `KEY` - Pull value from current environment

These are available during the build process via `ARG` + `ENV` in the builder stage. Use for frameworks that bake config at build time:
- Next.js `NEXT_PUBLIC_*` vars
- Vite `VITE_*` vars
- SvelteKit `$env/static/*`
- Any `process.env` accessed during build

**Runtime environment variables** should be passed when running the container:
```bash
docker run -e DATABASE_URL=postgres://... -e API_KEY=... myapp:latest
```

Use runtime env vars for:
- Secrets that shouldn't be baked into the image
- Config that changes per environment (dev/staging/prod)
- Frameworks that read env at runtime (SvelteKit `$env/dynamic/*`, Express `process.env`, etc.)

### Caching

Dockerfiles use BuildKit cache mounts for faster rebuilds:
- `# syntax=docker/dockerfile:1` - Enables BuildKit features
- `--mount=type=cache,target=...` - Caches package manager and build artifacts

#### Install Phase Caches

| Package Manager | Cache Directory |
|-----------------|-----------------|
| npm | `/root/.npm` |
| yarn v1 | `/usr/local/share/.cache/yarn` |
| yarn berry | `/root/.yarn/berry/cache` |
| pnpm | `/root/.local/share/pnpm/store` |
| bun | `/root/.bun/install/cache` |

Additional install caches:
- **Cypress**: `/root/.cache/Cypress` (if `cypress` dependency detected)

#### Build Phase Caches

Framework-specific build caches for faster incremental builds:

| Framework | Cache Directory |
|-----------|-----------------|
| Next.js | `.next/cache` |
| Remix / React Router | `.cache`, `.react-router` |
| Vite | `node_modules/.vite` |
| TanStack Start | `node_modules/.vite` |
| Astro | `node_modules/.astro` |
| Nuxt | `node_modules/.cache` |
| All frameworks | `node_modules/.cache` (webpack, babel, eslint, etc.) |

Additional build caches:
- **Moon repo**: `.moon/cache` (if `.moon/workspace.yml` detected)

#### Custom Cache Directories

You can specify custom cache directories in `package.json`:

```json
{
  "cacheDirectories": [
    ".cache",
    "tmp/build-cache"
  ]
}
```

Each directory will be cached between builds using BuildKit cache mounts.

## Project Structure

```
coolpack/
├── main.go                          # Entry point
├── build.sh                         # Build script
├── .goreleaser.yml                  # GoReleaser config for releases
├── .github/workflows/
│   └── release.yml                  # GitHub Actions release workflow
├── cmd/coolpack/
│   ├── root.go                      # Root CLI command
│   ├── plan.go                      # Plan subcommand
│   ├── prepare.go                   # Prepare subcommand (Dockerfile generation)
│   ├── build.go                     # Build subcommand
│   ├── run.go                       # Run subcommand
│   └── version.go                   # Version subcommand
└── pkg/
    ├── app/
    │   ├── context.go               # App context (path, env, file helpers)
    │   └── plan.go                  # Plan struct
    ├── detector/
    │   ├── detector.go              # Main detector, registers providers
    │   └── types.go                 # Provider interface
    ├── generator/
    │   └── generator.go             # Dockerfile generation
    ├── version/
    │   └── version.go               # Version info and update checker
    └── providers/node/
        ├── node.go                  # Node.js provider
        ├── package_json.go          # package.json parsing
        ├── package_manager.go       # Package manager detection
        ├── version.go               # Node version detection
        ├── framework.go             # Framework detection
        ├── config_parser.go         # JS/TS config parsing (tree-sitter)
        └── native_deps.go           # Native dependency detection
```

## Config File Parsing

Uses [tree-sitter](https://github.com/smacker/go-tree-sitter) for parsing JavaScript/TypeScript config files:
- `next.config.ts/js/mjs` - Detects `output: 'export'` for static builds
- `react-router.config.ts/js` - Detects `ssr: false` for SPA mode
- `nuxt.config.ts/js/mjs` - Detects `ssr: false` for SPA mode
- `astro.config.ts/js/mjs` - Detects `output: 'server'/'hybrid'` for SSR (default is static)
- `app.config.ts/js` - Detects `ssr: false` for Solid Start SPA mode, `server.preset: 'static'` for TanStack Start

## Dependencies

- `github.com/spf13/cobra` - CLI framework
- `github.com/smacker/go-tree-sitter` - AST parsing for JS/TS files

## Adding New Providers

1. Create `pkg/providers/<name>/<name>.go`
2. Implement `Provider` interface:
   - `Name() string`
   - `Detect(ctx *app.Context) (bool, error)`
   - `Plan(ctx *app.Context) (*app.Plan, error)`
3. Register in `pkg/detector/detector.go` `registerProviders()`

## Releases

Releases are automated via GitHub Actions using [GoReleaser](https://goreleaser.com/).

**To create a release:**
1. Create a new release on GitHub with a tag (e.g., `v1.0.0`)
2. GitHub Actions will automatically build binaries for:
   - Linux (amd64, arm64)
   - macOS (amd64, arm64)
   - Windows (amd64, arm64)
3. Binaries are attached to the release with checksums

**Files:**
- `.github/workflows/release.yml` - GitHub Actions workflow
- `.goreleaser.yml` - GoReleaser configuration
- `cmd/coolpack/version.go` - Version command

**Local build with version:**
```bash
VERSION=v1.0.0 ./build.sh
```
