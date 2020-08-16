FROM levonet/golang:go2go

RUN apt-get update && apt-get install git -y

RUN git clone https://github.com/dgrijalva/lfu-go /go/src/github.com/dgrijalva/lfu-go && \
    git clone https://github.com/arschles/assert /go/src/github.com/arschles/assert

COPY ./ /go/src/myapp

WORKDIR /go/src/myapp