#!/bin/bash
# Install Go programming language

echo "Installing Go..."
sudo apt update && sudo apt install -y golang-go

echo "Verifying installation..."
go version

echo "Done! Go is installed."
