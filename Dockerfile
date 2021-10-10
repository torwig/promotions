FROM golang:1.14 as builder

RUN mkdir /build
COPY . /build
WORKDIR /build/cmd/promotions
RUN CGO_ENABLED=0 go build -o /app .


FROM alpine:3.9

COPY --from=builder /app /app
COPY --from=builder /build/cert/localhost.key localhost.key
COPY --from=builder /build/cert/localhost.crt localhost.crt
RUN apk add --no-cache ca-certificates
CMD ["/app"]
