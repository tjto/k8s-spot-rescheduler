ARG VERSION=undefined

FROM golang:1.12 AS builder
ARG VERSION

RUN mkdir /k8s-spot-rescheduler

WORKDIR /k8s-spot-rescheduler
COPY *.go ./
COPY deploy deploy/
COPY metrics metrics/
COPY nodes nodes/
COPY scaler scaler/
COPY go.mod .
COPY go.sum .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-X main.VERSION=${VERSION}" -a -o k8s-spot-rescheduler

FROM alpine:3.9
RUN apk --no-cache add ca-certificates
WORKDIR /bin
COPY --from=builder /k8s-spot-rescheduler/k8s-spot-rescheduler .

ENTRYPOINT ["/bin/k8s-spot-rescheduler"]
