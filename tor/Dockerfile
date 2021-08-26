FROM alpine:latest

RUN apk update && apk add tor curl

RUN chown -R tor /etc/tor /var/lib/tor
RUN chmod -R 700 /var/lib/tor

USER tor
CMD tor -f /etc/tor/torrc