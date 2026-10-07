FROM golang:1.26 AS base
WORKDIR /app/api

FROM base AS deps
COPY api/go.mod api/go.sum ./
RUN go mod download

FROM deps AS build
COPY api/ ./
RUN CGO_ENABLED=0 GOOS=linux go build -o /southpark-messaging-api "./cmd"

FROM alpine:3.24 AS runtime
COPY --from=build /southpark-messaging-api /southpark-messaging-api
EXPOSE 3000

CMD ["/southpark-messaging-api"]
