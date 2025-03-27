FROM golang:1.24

# Set destination for COPY
WORKDIR /app
 
# Copies everything from your root directory into /app
COPY . .
 
# Installs Go dependencies
RUN go mod download
 
# Specifies the executable command that runs when the container starts
CMD ["go", "run", "main.go"]