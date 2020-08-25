# Build image
FROM golang:1.14.6 AS build-env

WORKDIR /app
ENV GO111MODULE=on
ENV CGO_ENABLED=0
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go test ./...
RUN go build

# Run image
FROM alpine:3.10.2
RUN apk add --no-cache ca-certificates
COPY --from=build-env /app/flyte-jira .

ENTRYPOINT ["./flyte-jira"]
