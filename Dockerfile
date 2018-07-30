FROM alpine:3.8
# Setup distro and user
RUN apk update && apk upgrade
RUN apk add go git ca-certificates
RUN adduser -h /opt/orangeforum -g 'orangeforum,,,,' -D orangeforum
RUN mkdir -p /opt/orangeforum

# Build orangeforum
COPY . /usr/src/orangeforum
WORKDIR /usr/src/orangeforum
RUN chown orangeforum -R /opt/orangeforum/
RUN go get -u github.com/s-gv/orangeforum/
RUN go build
RUN cp orangeforum /usr/bin/orangeforum
RUN apk del go git

# Setup and run orangeforum
USER orangeforum
WORKDIR /opt/orangeforum
RUN ./orangeforum -migrate
RUN ./orangeforum -createsuperuser
VOLUME /opt/orangeforum
CMD ./orangeforum
