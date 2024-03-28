FROM golang:1.22-alpine as build
WORKDIR /build
COPY app.go go.mod go.sum /build/
RUN go get -d -v ./...
RUN go build -ldflags="-s -w"

FROM gcr.io/distroless/static-debian12
WORKDIR /app
COPY --from=build /build/speisekarte /app
COPY index.html /app
CMD ["/app/speisekarte"]
