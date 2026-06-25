# Changelog

## Unreleased

- Renamed the app from WebChat to TinyChat.
- Refactored the project into `src/` Go files and `src/templates/index.html`.
- Added Docker Compose as the primary local development workflow.
- Switched the runtime image to Alpine with standard public CA certificates.
- Removed the bundled Ollama service; AI endpoints are external/user-provided.
- Added runtime OpenAI-compatible endpoint modal, one-hour browser local storage persistence, and testing-only warning copy.
- Added `/api/health`, `/api/models`, and `/api/chat` proxying with outbound debug logs.
- Added streaming chat rendering, live thinking display, optional thinking toggle, token/timing stats, and stream flushing.
- Added Markdown rendering for headings, lists, rules, blockquotes, inline code, fenced code blocks, and code copy buttons.
- Added an embedded SVG favicon.
- Added GitHub Actions workflows for Go formatting/vet/tests and Trivy Dockerfile/config scanning.
- Added a Docker release workflow that publishes TinyChat images to GHCR on GitHub Release publication.
- Added README attribution to the `llama.cpp` project that inspired TinyChat.
- Hardened the Docker runtime and proxy by removing baked-in env defaults, running as an unprivileged user, adding browser security headers, and filtering unsafe upstream response headers.
- Updated response stats to show icon-led processing speed, generation speed, total tokens, and total time.
- Updated live thinking so it sticks above the response while active and automatically collapses when generation begins.
- Moved the thinking toggle above generated answers and hid assistant copy actions until generation completes.
- Changed assistant copy to a visible footer icon that appears with response stats.
- Updated the model label to show context size and collapse file-path model IDs to filenames.
- Added auto-growing prompt input capped at half the viewport height.
- Centered the welcome state and composer on the empty start screen, then animated the composer to the bottom after the first prompt.
- Refreshed the UI theme with a modern dark teal palette, warmer accents, and clearer focus states.
- Added the UI screenshot to the README.
- Restored temporary browser-provided endpoints for local testing and documented `OPENAI_CHAT_HOST` as the production/shared deployment path.
- Reordered the README around description, server deployment, and environment variables.
- Simplified README hosting instructions into Native, Docker, and Docker Compose options.
- Removed the separate private/self-signed CA example from the README hosting section.
- Improved composer and message UI, including bottom alignment, message copy actions, model label with token context, and scroll behavior for long histories.
