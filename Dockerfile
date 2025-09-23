FROM golang:1.23 as builder
ARG APP_VERSION
ENV VERSION=${APP_VERSION}
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-X 'main.Version=${VERSION}'" -o ./out/server ./cmd/server/main.go

# Merge all compiled and make as a single simple docker image
FROM alpine:3.18 as server
MAINTAINER "OnePlatform Development"
WORKDIR /app
COPY --from=builder /app/out /app
RUN mkdir -p /app/storage

# Install tzdata and set docker image timezone to bangkok timezone
RUN apk add tzdata
ENV TZ=Asia/Bangkok

# Show app version
RUN /app/server version

EXPOSE 8000
CMD ["/app/server", "start"]
