FROM golang:1.14

# INSTALL BEE
ARG WORKDIR=/usr/src/app
RUN mkdir -p $WORKDIR
WORKDIR $WORKDIR
RUN go mod init beeinstall.org
RUN go get -v github.com/beego/bee/...

# GENERATE PROJECT
ARG PROJECT=bee_sample
RUN bee new $PROJECT
WORKDIR $GOPATH/src/$PROJECT
RUN go get -v

EXPOSE 8080

CMD ["bee", "run"]
