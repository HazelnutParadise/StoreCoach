FROM oven/bun:latest AS bun
WORKDIR /app
COPY . .
RUN cd frontend && bun install --frozen-lockfile && bun run build

FROM golang:1.23-alpine
WORKDIR /app
COPY --from=bun /app/frontend/dist /frontend/dist
COPY . .
RUN go build -o main .
RUN chmod +x main

CMD ["./main"]