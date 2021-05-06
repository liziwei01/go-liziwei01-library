#version 1.0
FROM golang:latest
LABEL maintainer="alssylk@gmail.com"
RUN ["mkdir", "-p", "/home/work/go-liziwei01-library"]
WORKDIR /home/work/go-liziwei01-library
COPY . /home/work/go-liziwei01-library
CMD ["/home/work/go-liziwei01-library/docker_run"] 
