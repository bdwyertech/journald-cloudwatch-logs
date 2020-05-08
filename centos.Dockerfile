FROM centos:7

MAINTAINER Brian Dwyer

RUN curl -L https://dl.google.com/go/go1.14.2.linux-amd64.tar.gz | tar xzf - --directory /usr/local

ENV PATH="/root/go/bin:/usr/local/go/bin:${PATH}"

RUN yum install -y openssl-devel systemd-devel gcc git \
    && go get honnef.co/go/tools/cmd/staticcheck

WORKDIR /go/src/journald-cloudwatch-logs

COPY . .

RUN go build -v .

RUN mkdir /var/lib/journald-cloudwatch-logs/
