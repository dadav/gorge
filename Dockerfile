FROM alpine:3.21@sha256:21dc6063fd678b478f57c0e13f47560d0ea4eeba26dfc947b2a4f81f686b9f45
RUN adduser -k /dev/null -u 10001 -D gorge \
  && chgrp 0 /home/gorge \
  && chmod -R g+rwX /home/gorge
COPY gorge /
USER 10001
VOLUME [ "/home/gorge" ]
ENTRYPOINT ["/gorge"]
CMD [ "serve" ]
