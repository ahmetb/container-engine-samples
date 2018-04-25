FROM golang:1.8-alpine
ADD . /go/src/hello-app
RUN apk add --no-cache git && \
	go get github.com/prometheus/client_golang/prometheus/...
RUN go install hello-app

FROM alpine:latest
COPY --from=0 /go/bin/hello-app .
ENV PORT 8080
CMD ["./hello-app"]
