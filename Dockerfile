FROM golang:1.12 as builder
WORKDIR /
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix .

FROM alpine:latest
WORKDIR /app/
COPY --from=builder /olivsoft-golang-api /app/olivsoft-api
EXPOSE 8888
ENTRYPOINT ./olivsoft-api