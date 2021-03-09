FROM golang:1.14.2

RUN mkdir /app

ADD . /app

WORKDIR /app

RUN go build -o main .

#RUN /bin/bash -c "source ./env"
#RUN echo $MGDBURL

CMD ["/app/main"]
