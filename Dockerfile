FROM golang:1.21-alpine as build
WORKDIR /build
COPY app.go go.mod go.sum /build/
RUN go get -d -v ./...
RUN go build -ldflags="-s -w"

FROM alpine:3.18
WORKDIR /app
COPY --from=build /build/speisekarte /app
COPY index.html /app
CMD ["/app/speisekarte"]
