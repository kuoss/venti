FROM alpine:3.15
ARG VENTI_VERSION
COPY --from=venti:base3 /app /app
WORKDIR /app
ENTRYPOINT ["/app/venti"]
