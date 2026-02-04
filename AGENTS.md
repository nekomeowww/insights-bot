# Repository Guidelines

## Project Structure & Module Organization
- `cmd/insights-bot/` contains the main entrypoint for the bot binary.
- `internal/` hosts core application logic (bots, services, configs, datastore).
- `pkg/` provides reusable packages (i18n, logging, utils, health checks).
- `ent/` is generated ORM code; schemas live in `ent/schema/`.
- `locales/` holds translation files (e.g., `en.yaml`, `zh-CN.yaml`).
- `docs/` contains documentation assets; `production/` includes deployment helpers.

## Build, Test, and Development Commands
- `go build -a -o build/insights-bot github.com/nekomeowww/insights-bot/cmd/insights-bot` builds a local binary.
- `go build -a -o release/insights-bot github.com/nekomeowww/insights-bot/cmd/insights-bot` creates a release artifact.
- `docker compose --profile hub up -d` runs the prebuilt image.
- `docker compose --profile local up -d --build` builds and runs from local code.
- `docker buildx build --platform linux/arm64,linux/amd64 -t <tag> -f Dockerfile .` builds multi-arch images.

## Coding Style & Naming Conventions
- Language: Go; use `gofmt` on all Go files.
- Avoid hand-editing `ent/` generated files; update schemas and re-run codegen.
- Prefer descriptive package names and keep exported identifiers in `CamelCase`.

## Testing Guidelines
- Tests are Go `*_test.go` files (e.g., `internal/models/chathistories/recap_test.go`).
- Run all tests with `go test ./...`.
- If you add new logic, add or extend tests in the nearest package.

## Commit & Pull Request Guidelines
- Commit messages follow Conventional Commits: `feat:`, `fix:`, `chore:` with optional scopes (e.g., `chore(deps): ...`). `[i18n] ...` is also used for translation updates.
- PRs should include a concise summary, testing notes (`go test ./...`, docker profile used), and any config changes.
- For behavior changes, mention impacted bots or commands (e.g., `/recap_forwarded`).

## Security & Configuration Tips
- Copy `.env.example` to `.env` for local runs; never commit secrets.
- Key ports: 6060 (pprof), 7069 (health), 7070–7072 (bot webhooks).
- Regenerate ORM after schema changes: `go generate ./ent`.
