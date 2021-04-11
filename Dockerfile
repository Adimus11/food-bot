FROM golang:1.13

RUN mkdir app
WORKDIR /app

COPY . .

RUN go get
RUN go build *.go

ENTRYPOINT ["./main"]
