FROM ubuntu:bionic
ENV args=""
# Setup distro and user
RUN apt-get update && apt-get upgrade -y
RUN apt-get install -y golang git ca-certificates gcc sqlite make
RUN mkdir -p /opt/orangeforum
RUN adduser --home /opt/orangeforum --gecos 'orangeforum,,,,' --disabled-password orangeforum

# Build orangeforum from source
COPY . /usr/src/orangeforum
WORKDIR /usr/src/orangeforum
RUN chown orangeforum -R /opt/orangeforum/
RUN go get -u github.com/s-gv/orangeforum/
RUN go build
RUN cp orangeforum /usr/bin/orangeforum

# Cleanup build and dependencies
RUN apt-get purge -y golang git

# Setup and run orangeforum
USER orangeforum
WORKDIR /opt/orangeforum
RUN orangeforum -migrate
VOLUME /opt/orangeforum
CMD orangeforum $args
