FROM alpine:latest

RUN mkdir -p /app

COPY heartbeatApp /app

CMD [ "/app/heartbeatApp"]
