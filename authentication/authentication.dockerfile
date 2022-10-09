FROM alpine:latest

RUN mkdir -p /app

COPY authenticationApp /app

CMD [ "/app/authenticationApp"]
