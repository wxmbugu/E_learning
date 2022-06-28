# syntax=docker/dockerfile:1
FROM golang:1.18
WORKDIR /go/src/github.com/worker
COPY go.mod go.sum main.go  ./
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .
FROM jrottenberg/ffmpeg:4-alpine


FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
ENV LD_LIBRARY_PATH=/usr/local/lib
COPY --from=jrottenberg/ffmpeg / /
COPY app.env /root
COPY --from=0 /go/src/github.com/worker/app .
CMD ["./app"]


