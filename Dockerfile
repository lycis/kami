FROM golang:1.9.2

LABEL maintainer="daniel@deder.at"
LABEL licence="GPLv3"

EXPOSE 8080
VOLUME /usr/lib/kami

WORKDIR /go/src/gitlab.com/lycis/kami
COPY . .

RUN curl https://glide.sh/get | sh
RUN glide update
RUN go test -v -cover gitlab.com/lycis/kami/...
RUN go install gitlab.com/lycis/kami

ENTRYPOINT ["/go/bin/kami", "--lib", "/usr/lib/kami/"]
