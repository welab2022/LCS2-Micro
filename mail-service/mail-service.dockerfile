FROM alpine:latest

RUN mkdir -p /app

COPY mailerApp /app
COPY templates /templates

CMD [ "/app/mailerApp"]
