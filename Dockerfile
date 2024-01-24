# Make small image for running this service
FROM  ubuntu:latest

LABEL maintainer="manhtokim@gmail.com"

RUN apt-get update \
    && apt-get install -y \
        ca-certificates \
        git \
        gcc \
        g++ \
        libc-dev \
        bash
        
WORKDIR /user/loacl/bin
COPY --from=intern-backend-builder:builder /projects/phenikaa_intern/phenikaa_intern_be/. /usr/local/bin/

RUN ls -ls /usr/local/bin/infrastructure/

CMD ["./phenikaa_intern_be"]