FROM golang:1.22-alpine AS build

ARG GOPROXY=https://goproxy.cn,direct
ENV GOPROXY=${GOPROXY}
WORKDIR /src
COPY go.mod go.sum* ./
RUN go mod download all
COPY . .
RUN go build -mod=mod -o /bin/eip-api ./cmd/server

FROM alpine:3.20

COPY --from=build /bin/eip-api /bin/eip-api
COPY --from=build /src/internal/db/migrations /migrations
EXPOSE 8080
CMD ["/bin/eip-api"]
