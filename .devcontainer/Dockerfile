# Start from the official Golang image
FROM golang:1.21-bullseye

# Install sudo as root
RUN apt-get update && apt-get install -y sudo

# Create a new vscode user with sudo access
RUN adduser --disabled-password --gecos '' vscode && \
    adduser vscode sudo && \
    echo '%sudo ALL=(ALL) NOPASSWD:ALL' >> /etc/sudoers

# Use user vscode from now on
USER vscode

# Install Redis and Ngrok
RUN sudo apt-get install -y redis-server wget unzip

# Download and unzip ngrok
RUN curl -JLO https://bin.equinox.io/c/bNyj1mQVY4c/ngrok-v3-stable-linux-amd64.tgz && \
    sudo tar xvzf ngrok-v3-stable-linux-amd64.tgz -C /usr/local/bin

RUN go install -v golang.org/x/tools/gopls@latest && go install -v github.com/go-delve/delve/cmd/dlv@latest
