# Build image
FROM golang:1.12 AS build-env

WORKDIR /app
ENV GO111MODULE=on
ENV CGO_ENABLED=0
COPY go.mod go.sum ./
COPY . .
RUN go test ./...
RUN go build

# Run image
FROM alpine:3.7
COPY --from=build-env /app/flyte-jira .


ENTRYPOINT ["./flyte-jira"]
