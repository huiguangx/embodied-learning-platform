FROM golang:1.22-alpine AS build

WORKDIR /src
COPY go.mod go.sum* ./
RUN go mod download
COPY . .
RUN go build -mod=mod -o /bin/eip-api ./cmd/server

FROM alpine:3.20

COPY --from=build /bin/eip-api /bin/eip-api
COPY --from=build /src/internal/db/migrations /migrations
EXPOSE 8080
CMD ["/bin/eip-api"]
