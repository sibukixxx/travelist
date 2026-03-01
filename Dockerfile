# Production: single container (Go serves frontend static files)

# Stage 1: Build frontend
FROM node:22-alpine AS frontend-builder

WORKDIR /app/frontend
COPY frontend/package.json frontend/package-lock.json ./
RUN npm ci
COPY frontend/ .
RUN npm run build

# Stage 2: Build Go API
FROM golang:1.24-alpine AS api-builder

WORKDIR /app/api
COPY api/go.mod api/go.sum ./
RUN go mod download
COPY api/ .
RUN CGO_ENABLED=0 go build -o /server ./cmd/server

# Stage 3: Runtime
FROM alpine:3.21

RUN apk add --no-cache ca-certificates tzdata

WORKDIR /app
COPY --from=api-builder /server /app/server
COPY --from=frontend-builder /app/frontend/dist /app/dist
COPY api/prompts/ /app/prompts/

EXPOSE 8080

ENV STATIC_DIR=/app/dist

CMD ["/app/server"]
