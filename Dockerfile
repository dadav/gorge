FROM alpine:3.20
RUN adduser -k /dev/null -u 10001 -D gorge \
  && chgrp 0 /home/gorge \
  && chmod -R g+rwX /home/gorge
COPY gorge /
USER 10001
VOLUME [ "/home/gorge" ]
ENTRYPOINT ["/gorge"]
CMD [ "serve" ]
