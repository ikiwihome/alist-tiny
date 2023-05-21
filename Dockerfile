FROM alpine:latest as builder
LABEL stage=go-builder
WORKDIR /app/
COPY ./ ./
RUN apk add --no-cache bash curl gcc git go musl-dev; \
    bash build.sh release docker

FROM alpine:latest
LABEL MAINTAINER="ikiwicc@mail.com"
WORKDIR /opt/alist/
COPY --from=builder /app/bin/alist ./
COPY entrypoint.sh /entrypoint.sh
RUN apk add --no-cache bash nano wget zip ca-certificates su-exec tzdata
RUN chmod +x /entrypoint.sh
ENV PUID=0 PGID=0 UMASK=022
EXPOSE 5244
CMD [ "/entrypoint.sh" ]