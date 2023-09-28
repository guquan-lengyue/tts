FROM golang:alpine as build
WORKDIR /src/
COPY . /src/
ENV GO111MODULE=on
ENV GOPROXY=https://goproxy.cn
RUN go build -o app.bin

FROM alpine
WORKDIR /app
ENV GIN_MODE=release
COPY --from=build /src/app.bin .
ENTRYPOINT ["./app.bin"]
EXPOSE 2580

