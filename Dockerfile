FROM golang:alpine

RUN apk update && apk add git && rm -rf /var/cache/apk/*

ADD . /go/src/code.olipicus.com/trueselect_coupon
WORKDIR go/src/code.olipicus.com/trueselect_coupon

RUN go get github.com/line/line-bot-sdk-go/linebot
RUN go build -o line .
CMD ["./line"]
