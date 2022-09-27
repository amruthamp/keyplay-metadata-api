FROM alpine:3.9

RUN apk --no-cache --update add ca-certificates

#copy swagger ui static files
COPY swaggerui /swaggerui

# copy generated swagger spec
ADD ./swagger.json /swaggerui

ADD ./build /go/bin

EXPOSE 8080
