FROM oven/bun:latest AS frontend
WORKDIR /app
COPY . .
RUN cd frontend && bun install --frozen-lockfile && bun run build

FROM golang:1.23-alpine AS backend
WORKDIR /app
COPY --from=frontend /app/frontend/dist /frontend/dist
COPY . .
RUN CGO_ENABLED=0 go build -o main .

FROM alpine:latest

COPY --from=backend /app/main /main
COPY --from=backend /app/StoreCoach.env /StoreCoach.env

CMD ["/main"]