## Builder Image
FROM golang:1.19 as builder
ENV GO111MODULE=on
COPY . /otel
WORKDIR /otel
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o bin/application

## Run Image
FROM alpine:3.16
COPY --from=builder /otel /otel
COPY --from=builder /otel/bin/application application
EXPOSE 1323
ENTRYPOINT ["./application"]