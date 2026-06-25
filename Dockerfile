FROM golang:1.22-alpine AS build

WORKDIR /src
COPY go.mod ./
COPY src ./src
RUN CGO_ENABLED=0 GOOS=linux go build -buildvcs=false -trimpath -ldflags="-s -w" -o /out/tinychat ./src

FROM alpine:3.21
RUN apk add --no-cache ca-certificates
COPY --from=build /out/tinychat /tinychat
ENV APP_HOST=0.0.0.0
ENV APP_PORT=8080
ENV OPENAI_CHAT_HOST=
EXPOSE 8080
ENTRYPOINT ["/tinychat"]
