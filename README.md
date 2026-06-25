# TinyChat

Extremely lightweight Go web chat UI for OpenAI-compatible chat APIs, intended to feel close to the built-in `llama.cpp` web UI while staying tiny and dependency-free.

TinyChat was inspired by the built-in web UI from [`llama.cpp`](https://github.com/ggml-org/llama.cpp).

![TinyChat screenshot](screen.png)

## Host TinyChat

### Native

Download the binary for your platform from the latest GitHub Release.

#### Linux

```sh
curl -L -o tinychat-linux-amd64.tar.gz https://github.com/kozmixb/tinychat/releases/latest/download/tinychat-linux-amd64.tar.gz
tar -xzf tinychat-linux-amd64.tar.gz
sudo install -m 0755 tinychat /usr/local/bin/tinychat
```

Use `tinychat-linux-arm64.tar.gz` for ARM64 Linux hosts.

Create `/etc/systemd/system/tinychat.service`:

```ini
[Unit]
Description=TinyChat
After=network-online.target
Wants=network-online.target

[Service]
Environment=APP_HOST=0.0.0.0
Environment=APP_PORT=8080
Environment=OPENAI_CHAT_HOST=https://api.example.com/v1
ExecStart=/usr/local/bin/tinychat
Restart=on-failure
RestartSec=5
User=tinychat
Group=tinychat
NoNewPrivileges=true

[Install]
WantedBy=multi-user.target
```

Create the service user and start TinyChat:

```sh
sudo useradd --system --no-create-home --shell /usr/sbin/nologin tinychat
sudo systemctl daemon-reload
sudo systemctl enable --now tinychat
```

#### macOS

```sh
curl -L -o tinychat-darwin-arm64.tar.gz https://github.com/kozmixb/tinychat/releases/latest/download/tinychat-darwin-arm64.tar.gz
tar -xzf tinychat-darwin-arm64.tar.gz
chmod +x tinychat
OPENAI_CHAT_HOST=https://api.example.com/v1 APP_HOST=0.0.0.0 APP_PORT=8080 ./tinychat
```

Use `tinychat-darwin-amd64.tar.gz` for Intel Macs.

#### Windows

```powershell
Invoke-WebRequest -Uri "https://github.com/kozmixb/tinychat/releases/latest/download/tinychat-windows-amd64.zip" -OutFile "tinychat-windows-amd64.zip"
Expand-Archive .\tinychat-windows-amd64.zip -DestinationPath .\tinychat
$env:OPENAI_CHAT_HOST="https://api.example.com/v1"
$env:APP_HOST="0.0.0.0"
$env:APP_PORT="8080"
.\tinychat\tinychat.exe
```

Open `http://your-server:8080` after starting TinyChat.

### Docker

Run the container and configure the OpenAI-compatible endpoint through the container environment:

```sh
docker run --rm -p 8080:8080 -e APP_HOST=0.0.0.0 -e OPENAI_CHAT_HOST=https://api.example.com/v1 tinychat
```

Open `http://your-server:8080`.

### Docker Compose

```sh
OPENAI_CHAT_HOST=https://api.example.com/v1 docker compose up -d --build
```

Open `http://your-server:8080`.

Temporary browser-provided endpoints are intended only for local testing. For production or shared deployments, set `OPENAI_CHAT_HOST` on the host server.

## Build With Docker

```sh
docker build -t tinychat .
```

## Environment Variables

| Variable           | Default     | Description                                                                                                                                                                                    |
| ------------------ | ----------- | ---------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------- |
| `APP_PORT`         | `8080`      | HTTP port for the web UI.                                                                                                                                                                      |
| `APP_HOST`         | `127.0.0.1` | Bind host for standalone runs. Docker sets this to `0.0.0.0`.                                                                                                                                  |
| `OPENAI_CHAT_HOST` | empty       | OpenAI-compatible API base URL. `http://` is added when the scheme is omitted; use `https://...` for HTTPS endpoints. If empty, the UI can accept a temporary endpoint for local testing only. |

## License

MIT. See [LICENSE](LICENSE).
