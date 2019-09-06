# Dockerfile References: https://docs.docker.com/engine/reference/builder/

# Start from the latest golang base image
FROM golang:latest

# Add Maintainer Info
LABEL maintainer="Robert Rong <robert.rong@agile-solutions.ch>"

# Set the Current Working Directory inside the container
WORKDIR /app

# first GO build and then copy this into the workdir
COPY main .

# Expose port 8080 to the outside world
EXPOSE 8080

# Command to run the executable
ENTRYPOINT ["/bin/bash"]