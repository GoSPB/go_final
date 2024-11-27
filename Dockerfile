# base
FROM golang:1.23.2-alpine3.20 AS base

# build
FROM base AS build
WORKDIR /src
COPY . .
RUN go mod download
RUN go build -o todo

# app
FROM base
WORKDIR /app
COPY --from=build /src/todo .
COPY .env .
COPY web ./web/
CMD ["/app/todo"]
