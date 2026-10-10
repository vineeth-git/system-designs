FROM golang:1.27-alpine AS build
WORKDIR /src
COPY src/ ./
RUN go mod tidy && CGO_ENABLED=0 go build -o /shortener .


FROM alpine:3.20
COPY --from=build /shortener /shortener
EXPOSE 8000
ENTRYPOINT ["/shortener"]
