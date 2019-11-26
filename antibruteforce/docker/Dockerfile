FROM golang:1.13
RUN mkdir /code
ADD . /code
WORKDIR /code
RUN go mod vendor && go mod download