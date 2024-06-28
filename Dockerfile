# Stage 1: Install Terraform on an Ubuntu base image
FROM ubuntu:latest as terraform

# Install dependencies for Terraform
RUN apt-get update \
    && apt-get install -y wget unzip

# Download and install Terraform
RUN wget https://releases.hashicorp.com/terraform/0.15.4/terraform_0.15.4_linux_amd64.zip \
    && unzip terraform_0.15.4_linux_amd64.zip \
    && mv terraform /usr/local/bin/ \
    && rm terraform_0.15.4_linux_amd64.zip

# Stage 2: Build Go application using a Go base image
FROM golang:1.21.4-alpine as builder

WORKDIR /app

# Copy go.mod and go.sum files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go application
RUN CGO_ENABLED=0 GOOS=linux go build -o /tenant-management-api

# Stage 3: Create final runtime image
FROM alpine:latest

RUN wget https://dl.google.com/go/go1.21.4.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go1.21.4.linux-amd64.tar.gz && \
    rm go1.21.4.linux-amd64.tar.gz

# Set Go environment variables
ENV GOROOT=/usr/local/go
ENV GOPATH=/go
ENV PATH=$GOROOT/bin:$GOPATH/bin:$PATH

WORKDIR /root/

# Copy the Terraform binary from the first stage
COPY --from=terraform /usr/local/bin/terraform /usr/local/bin/

# Copy the built Go application from the builder stage
COPY --from=builder /tenant-management-api /tenant-management-api

# Expose the necessary port
EXPOSE 8000

# Set environment variables
ENV PORT 8000
ENV TF_EXECUTABLE "/usr/local/bin/terraform"
ENV TF_WORKDIR "/root/terraform/"
ENV MODULE_NAME "tenant_management"
ENV GOOGLE_CREDS_PATH "/creds.json"

# Command to run the application
CMD ["/tenant-management-api"]