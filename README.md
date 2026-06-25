# TinyChat

Extremely lightweight Go web chat UI for OpenAI-compatible chat APIs, intended to feel close to the built-in `llama.cpp` web UI while staying tiny and dependency-free.

If `OPENAI_CHAT_HOST` is not provided, the UI asks for an OpenAI-compatible API host at startup. HTTP and HTTPS endpoints are supported. Temporary browser-provided endpoints are remembered in a cookie for one hour; set the environment variable for repeatable deployments.

It proxies:

- `GET /api/models` to `${OPENAI_CHAT_HOST}/models`
- `POST /api/chat` to `${OPENAI_CHAT_HOST}/chat/completions`

## Configuration

| Variable | Default | Description |
| --- | --- | --- |
| `APP_PORT` | `8080` | HTTP port for the web UI. |
| `APP_HOST` | `127.0.0.1` | Bind host for standalone runs. Docker sets this to `0.0.0.0`. |
| `OPENAI_CHAT_HOST` | empty | OpenAI-compatible API base URL. `http://` is added when the scheme is omitted; use `https://...` for HTTPS endpoints. If empty, the UI asks for a temporary runtime value. |

## Local Development

```sh
docker compose up --build
```

Open `http://127.0.0.1:8080`.

Use Docker Compose for local development. It keeps the web UI runtime consistent across Windows and other development machines.

Compose does not start or depend on a local AI service. Provide an external OpenAI-compatible endpoint in the startup modal, or set it when starting compose:

```sh
OPENAI_CHAT_HOST=https://api.example.com/v1 docker compose up --build
```

## Run Standalone

```sh
go run -buildvcs=false ./src
```

Open `http://127.0.0.1:8080`.

The UI will ask for an API host unless `OPENAI_CHAT_HOST` is set.

For an external HTTPS endpoint:

```sh
OPENAI_CHAT_HOST=https://api.example.com/v1 go run -buildvcs=false ./src
```

On Windows PowerShell:

```powershell
$env:OPENAI_CHAT_HOST="https://api.example.com/v1"
go run -buildvcs=false ./src
```

## Run With Docker

```sh
docker build -t tinychat .
docker run --rm -p 8080:8080 -e APP_HOST=0.0.0.0 tinychat
```

To configure the API host through the container environment:

```sh
docker run --rm -p 8080:8080 -e APP_HOST=0.0.0.0 -e OPENAI_CHAT_HOST=https://api.example.com/v1 tinychat
```

The Docker image is based on Alpine and includes the standard public CA bundle. For private or self-signed CAs, mount an augmented CA bundle over the default path:

```sh
docker run --rm -p 8080:8080 -v /path/to/ca-bundle.crt:/etc/ssl/certs/ca-certificates.crt:ro tinychat
```

## License

MIT. See [LICENSE](LICENSE).
