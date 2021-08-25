FROM alpine:latest

RUN apk update && apk add tor iptables

COPY torrc /etc/tor/torrc
RUN chown -R tor /etc/tor /var/lib/tor

# Tor can't redirect packets to hosts other than localhost, so once redirect
# to localhost, then pass it to another host (api).
RUN systemctl enable iptables && iptables -t nat -A PREROUTING -p tcp --dport 8081 -j DNAT --to-destination api:8080

USER tor

CMD tor -f /etc/tor/torrc