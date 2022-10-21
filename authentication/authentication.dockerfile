FROM alpine:latest

# openssl is the only required thing to install
# RUN apk --update add openssl

RUN mkdir -p /app

COPY authenticationApp /app

CMD [ "/app/authenticationApp"]
