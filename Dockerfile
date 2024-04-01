FROM alpine
ADD build/webserver-api /
WORKDIR /app
CMD ["/webserver-api"]
