FROM centos:7

MAINTAINER Brian Dwyer

RUN curl -L https://dl.google.com/go/go1.14.2.linux-amd64.tar.gz | tar xzf - --directory /usr/local

ENV PATH="/root/go/bin:/usr/local/go/bin:${PATH}"

RUN yum install -y openssl-devel systemd-devel gcc git \
    && go get honnef.co/go/tools/cmd/staticcheck

WORKDIR /go/src/journald-cloudwatch-logs

COPY . .

ARG TEST
RUN if [ "x$TEST" != "x" ] ; then \
      go fmt $(go list ./... | grep -v /vendor/) | xargs -I {} -r /bin/sh -c "/bin/echo {} && exit 1" \
      && go vet $(go list ./... | grep -v /vendor/) \
      && staticcheck $(go list ./... | grep -v /vendor/) \
      && go test -v -race $(go list ./... | grep -v /vendor/); \
    fi

RUN go build -v .
