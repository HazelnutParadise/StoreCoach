FROM oven/bun:latest AS bun
WORKDIR /app
COPY . .
RUN cd frontend && bun install --frozen-lockfile && bun run build

FROM golang:1.23 AS build
WORKDIR /app
COPY . .
RUN go build -o /main

FROM scratch
WORKDIR /
COPY --from=build /app/main /main
COPY --from=bun /app/frontend/dist /frontend/dist

CMD ["/main"]