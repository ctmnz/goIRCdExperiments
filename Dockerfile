FROM golang

CMD mkdir -p /go/src/github.com/ctmnz/goirc

ADD . /go/src/github.com/ctmnz/goirc/

RUN go install github.com/ctmnz/goirc

ENTRYPOINT /go/bin/goirc

EXPOSE 6667


