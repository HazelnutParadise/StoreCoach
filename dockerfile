# Frontend build stage
FROM oven/bun:latest AS frontend
WORKDIR /app
# 只複製前端相關檔案
COPY . .
WORKDIR /app/frontend
RUN bun install --frozen-lockfile
RUN bun run build

# Backend build stage
FROM golang:1.23-alpine AS backend
WORKDIR /app
# 先複製前端構建結果
COPY --from=frontend /app ./

RUN CGO_ENABLED=0 go build -o main .

# Final stage
FROM alpine:latest
WORKDIR /app
COPY --from=backend /app/main ./
# COPY --from=backend /app/StoreCoach.env ./
COPY --from=frontend /app/frontend/dist ./frontend/dist

CMD ["./main"]