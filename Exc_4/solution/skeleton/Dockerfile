  # todo
FROM golang:1.27 AS bob_the_builder
COPY . /app
WORKDIR /app

RUN go install github.com/swaggo/swag/cmd/swag@latest
RUN scripts/build-openapi-docs.sh
RUN scripts/build-application.sh

FROM alpine:latest
COPY --from=bob_the_builder /app/ordersystem /app/ordersystem
WORKDIR /app

RUN chmod +x /app/ordersystem

EXPOSE 5555

CMD ["/app/ordersystem"]