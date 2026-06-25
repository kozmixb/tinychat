FROM golang:1.22-alpine AS build

WORKDIR /src
COPY go.mod ./
COPY src ./src
RUN CGO_ENABLED=0 GOOS=linux go build -buildvcs=false -trimpath -ldflags="-s -w" -o /out/tinychat ./src

FROM alpine:3.21
RUN apk add --no-cache ca-certificates
RUN addgroup -S tinychat && adduser -S -G tinychat tinychat
COPY --from=build /out/tinychat /tinychat
EXPOSE 8080
USER tinychat
ENTRYPOINT ["/tinychat"]
