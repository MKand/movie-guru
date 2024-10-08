FROM node:20-alpine

# Install necessary tools for building Go
RUN apk add --no-cache build-base wget ca-certificates

# Download and install Go 1.22.7 (you can change the version if needed)
RUN wget https://go.dev/dl/go1.22.7.linux-amd64.tar.gz && \
    tar -C /usr/local -xzf go1.22.7.linux-amd64.tar.gz && \
    rm go1.22.7.linux-amd64.tar.gz

# Add Go to the PATH
ENV PATH="/usr/local/go/bin:$PATH"
# Install Genkit and typescript
RUN npm install -g genkit
RUN npm install -g typescript

# Create a project directory
RUN mkdir /app

# Set the working directory
WORKDIR /app

# Expose the Genkit port (if you plan to run the Genkit server)
EXPOSE 4000

# Add the Genkit path to the PATH environment variable (important!)
ENV PATH="/usr/local/bin:$PATH"

CMD ["sleep", "infinity"]