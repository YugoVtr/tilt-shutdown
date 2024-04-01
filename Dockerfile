FROM alpine
RUN apk add entr
RUN mkdir /app
WORKDIR /app
ADD build/webserver-api ./
CMD ls webserver-api | entr -rn ./webserver-api
