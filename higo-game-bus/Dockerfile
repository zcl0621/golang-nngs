FROM golang:1.18-alpine AS build-dist
ENV GOPROXY='https://goproxy.cn,direct'

WORKDIR /go/cache

ADD go.mod .
ADD go.sum .
RUN go mod download

WORKDIR /go/release

ADD . .
RUN go mod tidy
RUN GOOS=linux CGO_ENABLED=0 go build -ldflags="-s -w" -tags netcgo -installsuffix cgo -o /bin/app main.go

FROM alpine as prod

COPY --from=build-dist /bin/app /bin/app

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories

RUN apk add --no-cache -U  tzdata && \
    ln -sf /usr/share/zoneinfo/Asia/Shanghai  /etc/localtime && \
    echo "Asia/Shanghai" > /etc/timezone && \
    chmod +x /bin/app

ENV RUN_MODE='release'
CMD ["/bin/app"]
EXPOSE 8080
EXPOSE 8888
