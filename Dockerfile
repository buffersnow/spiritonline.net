FROM golang:1.25.0-trixie AS builder
LABEL authors="Serafine Karayan"
WORKDIR /app
COPY . .
ARG SERVICE_NAME=gsp
RUN /app/scripts/run-service.sh $SERVICE_NAME
RUN mv /app/bin/${SERVICE_NAME}.lxb /app/bin/app

FROM debian:trixie-slim
WORKDIR /app
COPY --from=builder /app/bin/app /app/app
COPY --from=builder /app/.env /app/.env
ENTRYPOINT ["/app/app"]