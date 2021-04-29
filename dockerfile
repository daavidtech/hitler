FROM golang:1.16.3

WORKDIR /usr/src

COPY ./quotes ./quotes
COPY ./meinkampf ./meinkampf

ADD ./go.mod ./go.mod
ADD ./go.sum ./go.sum
ADD ./main.go ./main.go

RUN go get -d -v ./...
RUN go install -v ./...
RUN go build

ENTRYPOINT ["./hitler"]