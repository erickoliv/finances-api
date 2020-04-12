FROM golang:1.13 as builder
WORKDIR /
COPY . ./
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix .

FROM alpine:latest
WORKDIR /app/
COPY --from=builder /finances-api /app/run
EXPOSE 80
ENTRYPOINT ./run