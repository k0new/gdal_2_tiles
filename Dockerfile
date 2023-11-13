FROM golang:1.20
# Enable CGO for C code emulation
ENV CGO_ENABLED=1
# Move to /app
WORKDIR /app
# Update & install pkgs
RUN apt-get update -y && \
    apt install wget build-essential libncursesw5-dev libssl-dev \
    libsqlite3-dev tk-dev libgdbm-dev libc6-dev libbz2-dev libffi-dev zlib1g-dev -y && \
    apt-get install gdal-bin -y && \
    apt-get install libgdal-dev -y
# Copy go.mod & go.sum
ADD ./meteocontext/grib-tiler/go.mod .
ADD ./meteocontext/grib-tiler/go.sum .
# Install deps
RUN go mod download
# Copy project
COPY ./meteocontext/grib-tiler .
# Execute go run cmd/main.go


#FROM golang:1.20-bullseye AS builder
#
## Set necessary environmet variables needed for our image
#ENV CGO_ENABLED=0
## Move to working directory /build
#WORKDIR /build
## Copy and download dependency using go mod
#COPY go.mod .
#COPY go.sum .
#RUN go mod download
#
## Copy the code into the container
#COPY . .
## Build the application
## go build -o [name] [path to file]
#RUN go build -o app cmd/main.go
#
## Move to /dist directory as the place for resulting binary folder
#WORKDIR /dist
#
## Copy binary from build to main folder
#RUN cp /build/app .
#
#############################
## STEP 2 build a small image
#############################
#FROM ubuntu:latest
#
## Update & install pkgs
#RUN apt-get update -y && \
#    apt install wget build-essential libncursesw5-dev libssl-dev \
#    libsqlite3-dev tk-dev libgdbm-dev libc6-dev libbz2-dev libffi-dev zlib1g-dev -y && \
#    apt-get install gdal-bin -y && \
#    apt-get install libgdal-dev -y
#
#COPY . .
#COPY --from=builder /dist/app /
## Copy the code into the container
#
## Command to run the executable
#ENTRYPOINT ["./app"]