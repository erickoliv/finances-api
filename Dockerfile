FROM ubuntu
COPY olivsoft-golang-api /bin/main
ENTRYPOINT /bin/main
# Service listens on port 8080.
EXPOSE 8080