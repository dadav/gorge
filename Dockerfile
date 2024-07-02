FROM ubuntu:latest as user_stage
RUN useradd -u 10001 gorge

FROM scratch
ENV HOME /home/gorge
USER 10001
COPY gorge /
COPY --from=user_stage /etc/passwd /etc/passwd
ENTRYPOINT ["/gorge"]
CMD ["serve"]
