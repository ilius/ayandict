FROM golang:1.21-bookworm

WORKDIR /app

COPY . /app

#RUN apt-get update

RUN go build -tags nogui

CMD ./ayandict
