FROM alpine

#RUN mkdir /app

ADD bin/myblog .

#WORKDIR /app

#RUN go build -o main .


#RUN /bin/bash -c "source ./env"
#RUN echo $MGDBURL
EXPOSE 8080

ENTRYPOINT  ["./myblog"]
