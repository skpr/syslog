FROM previousnext/golang:1.12 as build
ADD . /go/src/github.com/skpr/syslog
WORKDIR /go/src/github.com/skpr/syslog
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o /usr/local/bin/syslog github.com/skpr/syslog

FROM skpr/base:1.x
COPY --from=build /usr/local/bin/syslog /usr/local/bin/syslog
VOLUME /dev/log
CMD ["syslog"]
