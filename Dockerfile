FROM alpine
RUN apk --no-cache add ca-certificates
WORKDIR /
COPY tmlist /
CMD [ "/tmlist" ]