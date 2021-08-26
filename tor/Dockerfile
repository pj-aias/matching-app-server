FROM alpine:latest

RUN apk update && apk add tor iptables

COPY torrc /etc/tor/torrc
RUN chown -R tor /etc/tor /var/lib/tor
RUN chmod -R 700 /var/lib/tor

USER tor

CMD tor -f /etc/tor/torrc