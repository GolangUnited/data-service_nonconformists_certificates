FROM golang:1.19-alpine as build
COPY . /src
WORKDIR /src
RUN go mod tidy
RUN go build ./cmd/server/main.go; \
    mv main certificate-server

FROM alpine
WORKDIR /app
COPY --from=build /src/certificate-server ./certificate-server
COPY --from=build /src/.env ./.env
CMD ["./certificate-server"]