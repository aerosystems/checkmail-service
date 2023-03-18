FROM alpine:latest
RUN mkdir /app

COPY ./checkmail-service/checkmail-service.bin /app

# Run the server executable
CMD [ "/app/checkmail-service.bin" ]