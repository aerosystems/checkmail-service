FROM alpine:3.20.0
RUN mkdir /app
RUN mkdir /app/logs

COPY ./checkmail-service.bin /app

# Run the server executable
CMD [ "/app/checkmail-service.bin" ]