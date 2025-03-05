FROM oven/bun:latest AS bun
COPY . .
RUN cd frontend && bun install && bun build

FROM golang:1.23 AS build
COPY . .
RUN go build -o /app

FROM scratch
COPY --from=build /app /app
COPY --from=bun /frontend/dist /frontend/dist

CMD ["/app"]