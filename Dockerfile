FROM alpine:3.20@sha256:1e42bbe2508154c9126d48c2b8a75420c3544343bf86fd041fb7527e017a4b4a
RUN adduser -k /dev/null -u 10001 -D gorge \
  && chgrp 0 /home/gorge \
  && chmod -R g+rwX /home/gorge
COPY gorge /
USER 10001
VOLUME [ "/home/gorge" ]
ENTRYPOINT ["/gorge"]
CMD [ "serve" ]
