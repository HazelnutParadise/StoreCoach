FROM oven/bun:1.4.4 AS frontend
WORKDIR /app/frontend
COPY frontend/package.json frontend/bun.lock ./
RUN bun install --frozen-lockfile
COPY frontend/ ./
RUN bun run build

FROM golang:1.27.3-alpine AS backend
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY main.go ./
COPY app/ ./app/
COPY database/ ./database/
COPY routes/ ./routes/
COPY --from=frontend /app/frontend/dist ./frontend/dist
RUN CGO_ENABLED=0 go build -trimpath -o /storecoach .

FROM alpine:latest
RUN apk add --no-cache ca-certificates
WORKDIR /app
COPY --from=backend /storecoach ./storecoach
EXPOSE 8080
CMD ["./storecoach"]
