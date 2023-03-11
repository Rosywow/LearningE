#docker image rm -f chatapp
#docker rm -f chatApp_backend
#docker build -t learnbackend -f Dockerfile .
#docker run -p 1235:1235 --name LEB --network chatnet --restart=always -d learnbackend
FROM golang:alpine AS build

ENV GO111MODULE=on
RUN go env -w GOPROXY=https://proxy.golang.com.cn,direct
WORKDIR /app

COPY . .


RUN go mod download

COPY *.go ./
RUN go build -o /chatapp

FROM alpine:latest

WORKDIR /app

COPY --from=build /chatapp /app/chatapp

EXPOSE 1235

#这个必须要加，不然没有权限

ENTRYPOINT ["/app/chatapp"]


