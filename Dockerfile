# syntax=docker/dockerfile:1
FROM golang:1.18
WORKDIR /go/src/github.com/e-learning
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .



FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY app.env /root
COPY --from=0 /go/src/github.com/e-learning/app .
CMD ["./app"]


