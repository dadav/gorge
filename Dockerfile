FROM ubuntu:latest
RUN useradd -u 10001 gorge

FROM scratch
ENV HOME /home/gorge
USER 10001
COPY gorge /
COPY --from=0 /etc/passwd /etc/passwd
ENTRYPOINT ["/gorge"]
