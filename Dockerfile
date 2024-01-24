# Make small image for running this service
FROM  ubuntu:latest

LABEL maintainer="manhtokim@gmail.com"

# Add the service binary to the container
RUN apt-get update
RUN apt-get install -y ca-certificates

WORKDIR /user/loacl/bin
COPY --from=intern-backend-builder:builder /projects/phenikaa_intern/phenikaa_intern_be/. /usr/local/bin/

RUN ls -ls /usr/local/bin/infrastructure/

CMD ["./phenikaa_intern_be"]