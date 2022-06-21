FROM scratch
COPY --from=nats-streaming:0.24.6-alpine3.15 /usr/local/bin/nats-streaming-server /nats-streaming-server
EXPOSE ${natsStrPort1} ${natsStrPort2}
ENTRYPOINT ["/nats-streaming-server"]
CMD ["-m", "${natsStrPort2}"]