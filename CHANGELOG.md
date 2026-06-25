# Changelog

## Unreleased

- Renamed the app from WebChat to TinyChat.
- Refactored the project into `src/` Go files and `src/templates/index.html`.
- Added Docker Compose as the primary local development workflow.
- Switched the runtime image to Alpine with standard public CA certificates.
- Removed the bundled Ollama service; AI endpoints are external/user-provided.
- Added runtime OpenAI-compatible endpoint modal, one-hour browser cookie persistence, and session-only warning copy.
- Added `/api/health`, `/api/models`, and `/api/chat` proxying with outbound debug logs.
- Added streaming chat rendering, live thinking display, optional thinking toggle, token/timing stats, and stream flushing.
- Added Markdown rendering for headings, lists, rules, blockquotes, inline code, fenced code blocks, and code copy buttons.
- Added an embedded SVG favicon.
- Added GitHub Actions workflows for Go formatting/vet/tests and Trivy Dockerfile/config scanning.
- Added a Docker release workflow that publishes TinyChat images to GHCR on GitHub Release publication.
- Added README attribution to the `llama.cpp` project that inspired TinyChat.
- Hardened the Docker runtime and proxy by removing baked-in env defaults, running as an unprivileged user, adding browser security headers, and filtering unsafe upstream response headers.
- Improved composer and message UI, including bottom alignment, message copy actions, model label with token context, and scroll behavior for long histories.
