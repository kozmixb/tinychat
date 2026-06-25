# Repository Instructions

- Use Docker Compose for local development and verification. Do not use `go run` as the default local workflow.
- Do not start, add, or depend on a local Ollama service in Compose. The AI endpoint is external and provided by the user or client.
- Use an Alpine runtime image with standard public CA certificates. Users can mount an augmented CA bundle for private or self-signed CAs.
- Keep model reasoning/thinking hidden by default in the UI, but available through an explicit per-message toggle when the endpoint streams it.
- After code or configuration changes, rebuild and restart the app with Docker Compose so the running container reflects the latest files.
- Prefer `docker compose up -d --build` for the restart loop, then verify the served app over HTTP.
- Keep `CHANGELOG.md` updated with compact entries until the first release.
- Keep this file updated when workflow instructions change during development.
