FROM golang:1.17

RUN mkdir /app
ADD . /app
WORKDIR /app
RUN go build -o foodlogiq ./cmd/service
CMD ["/app/foodlogiq"]