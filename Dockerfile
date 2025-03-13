FROM alpine:3.21@sha256:a8560b36e8b8210634f77d9f7f9efd7ffa463e380b75e2e74aff4511df3ef88c

# Create non-root user and set up permissions in a single layer
RUN adduser -k /dev/null -u 10001 -D gorge \
  && chgrp 0 /home/gorge \
  && chmod -R g+rwX /home/gorge

# Copy application binary with explicit permissions
COPY --chmod=755 gorge /

# Set working directory
WORKDIR /home/gorge

# Switch to non-root user
USER 10001

# Define volume
VOLUME [ "/home/gorge" ]

# Set health check
HEALTHCHECK --interval=30s --timeout=3s \
  CMD curl -f http://localhost:8080/readyz || exit 1

ENTRYPOINT ["/gorge"]
CMD [ "serve" ]
