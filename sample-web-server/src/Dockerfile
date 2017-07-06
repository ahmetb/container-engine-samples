FROM golang:1.8-alpine
EXPOSE 8080
RUN apk add --no-cache git

WORKDIR /go/src/app
COPY . .
RUN go get -v -d ./...
RUN go build -o ./app
ENTRYPOINT ./app
