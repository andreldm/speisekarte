FROM golang:1.17-alpine as build
WORKDIR /build
COPY app.go go.mod go.sum /build/
RUN go get -d -v ./...
RUN go build

FROM alpine:3.13
WORKDIR /app
COPY --from=build /build/speisekarte /app
COPY index.html /app
CMD ["/app/speisekarte"]
