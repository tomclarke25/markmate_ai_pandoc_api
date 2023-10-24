FROM pandoc/core:latest

# Install Go
ENV GOLANG_VERSION 1.21.3
ENV CLERK_SECRET_KEY="dummy"
ENV AWS_LAMBDA_FUNCTION_VERSION="latest"
ENV AWS_LAMBDA_FUNCTION_NAME="testFunction"
ENV AWS_LAMBDA_FUNCTION_MEMORY_SIZE="128"
ENV _LAMBDA_SERVER_PORT=9000
ENV AWS_LAMBDA_RUNTIME_API=localhost:9000

RUN set -eux; \
    url="https://golang.org/dl/go${GOLANG_VERSION}.linux-amd64.tar.gz"; \
    wget -O go.tgz "$url"; \
    tar -C /usr/local -xzf go.tgz; \
    rm go.tgz;
ENV PATH $PATH:/usr/local/go/bin

# Set the working directory
WORKDIR /app

# Copy local code to the container image
COPY . .

# Build the Go app
RUN go build -o main

# Use the AWS Lambda runtime interface client as the entrypoint
ENTRYPOINT [ "/app/main" ]

EXPOSE 9000
