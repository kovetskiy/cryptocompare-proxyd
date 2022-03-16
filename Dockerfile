FROM alpine:edge

RUN apk --update --no-cache add ca-certificates

COPY /cryptocompare-proxyd /cryptocompare-proxyd

ENTRYPOINT ["/cryptocompare-proxyd"]
