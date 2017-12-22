FROM golang:1.9.2

LABEL maintainer="daniel@deder.at"
LABEL licence="GPLv3"

EXPOSE 23

WORKDIR /go/src/gitlab.com/lycis/kami/frontend/telnet
COPY . .

RUN curl https://glide.sh/get | sh
RUN glide update
RUN go test -v -cover gitlab.com/lycis/kami/frontend/telnet
RUN go install gitlab.com/lycis/kami/frontend/telnet

ENTRYPOINT ["/go/bin/telnet", "--port", "23", "--interface", "0.0.0.0"]
