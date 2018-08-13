FROM alpine:3.8
ENV args=""
# Setup distro and user
RUN apk update && apk upgrade
RUN apk add go git ca-certificates musl-dev sqlite
RUN mkdir -p /opt/orangeforum
RUN adduser -h /opt/orangeforum -g 'orangeforum,,,,' -D orangeforum

# Build orangeforum from source
COPY . /usr/src/orangeforum
WORKDIR /usr/src/orangeforum
RUN chown orangeforum -R /opt/orangeforum/
RUN go get -u github.com/s-gv/orangeforum/
RUN go build
RUN cp orangeforum /usr/bin/orangeforum

# Cleanup build and dependencies
RUN apk del go git musl-dev

# Setup and run orangeforum
USER orangeforum
WORKDIR /opt/orangeforum
RUN orangeforum -migrate
VOLUME /opt/orangeforum
CMD orangeforum $args
