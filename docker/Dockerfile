# Create a appuser in ubuntu
FROM ubuntu:latest AS builder
RUN useradd -u 10001 appuser

# Copy binary and appuser over to scratch image
FROM scratch
COPY bin/service /service
COPY --from=builder /etc/passwd /etc/passwd
USER appuser
EXPOSE 8081
ENTRYPOINT ["/service"]