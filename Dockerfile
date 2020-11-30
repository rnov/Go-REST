FROM golang:1 as builder

ENV APP=gorest

COPY . /${APP}/
WORKDIR /${APP}

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -mod=readonly -a -o /service/${APP} ./cmd/${APP}
RUN cp /${APP}/config/envs/docker/config.yml /service


FROM scratch
WORKDIR /

COPY --from=builder /service/* /
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

ENV ENV_PATH="config.yml"
EXPOSE 8080

CMD ["/gorest"]
