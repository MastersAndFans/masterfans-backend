FROM ubuntu:latest
LABEL authors="Lotus"

ENTRYPOINT ["top", "-b"]