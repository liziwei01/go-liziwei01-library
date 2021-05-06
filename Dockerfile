#version 1.0
FROM golang:latest
LABEL maintainer="alssylk@gmail.com"
RUN ["mkdir", "-p", "/home/work/github.com/go-liziwei01-library"]
WORKDIR /home/work/github.com/go-liziwei01-library
COPY . /home/work/github.com/go-liziwei01-library
CMD ["/home/work/github.com/go-liziwei01-library/docker_run"] 
