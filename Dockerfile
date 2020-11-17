FROM golang:1.15-alpine
WORKDIR /src 
COPY . /src
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o app .

FROM alpine:latest  
RUN \
    apk --no-cache add ca-certificates curl &&\
    mkdir -p /opt/trackip/config
WORKDIR /opt/trackip
COPY --from=0 /src/app .
CMD ["./app"]