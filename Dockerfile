#
#  TMList (c) 2023 by Mikhail Kondrashin (mkondrashin@gmail.com)
#  Copyright under MIT Lincese. Please see LICENSE file for details
#
#  Dockerfile - build docker image to run tmlist program
#

FROM alpine
RUN apk --no-cache add ca-certificates
WORKDIR /
COPY tmlist /
CMD [ "/tmlist" ]