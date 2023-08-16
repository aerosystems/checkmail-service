FROM alpine:latest
RUN mkdir /app
RUN mkdir /app/logs

COPY ./checkmail-service/checkmail-service.bin /app

# Run the server executable
CMD [ "/app/checkmail-service.bin" ]