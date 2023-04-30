FROM busybox:glibc
WORKDIR /
COPY tmlist /
CMD [ "/tmlist" ]
