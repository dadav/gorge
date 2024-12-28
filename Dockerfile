FROM alpine:3.21@sha256:21dc6063fd678b478f57c0e13f47560d0ea4eeba26dfc947b2a4f81f686b9f45

# Create non-root user and set up permissions in a single layer
RUN adduser -k /dev/null -u 10001 -D gorge \
  && chgrp 0 /home/gorge \
  && chmod -R g+rwX /home/gorge \
  # Add additional security hardening
  && chmod 755 /gorge

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
