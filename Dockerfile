FROM golang:1.23

# Install dependencies for OpenGL
RUN apt-get update && apt-get install -y libgl1-mesa-dev libxi-dev libxcursor-dev libxrandr-dev libxinerama-dev xorg-dev

# Set the working directory
WORKDIR /app

# Copy the project files
COPY . .

# Build the application
RUN go build -o hl7_to_fhir