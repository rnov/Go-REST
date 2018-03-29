FROM golang:1.8-alpine
# install needed packages
RUN apk update && \
 apk upgrade && \
 apk add git && \
 go get github.com/julienschmidt/httprouter && \
 go get github.com/go-redis/redis && \
 go get github.com/op/go-logging && \
 go get gopkg.in/yaml.v2
# config environment
RUN mkdir /go/src/Go-REST
ADD . /go/src/Go-REST
WORKDIR /go/src/Go-REST
# build (compile)
RUN cd /go/src/Go-REST/application &&  go build -o main .
# command to be launched
#CMD ["go", "run", "/go/src/Go-REST/application/main.go"]
