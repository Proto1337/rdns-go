# This Dockerfile uses two stages.
# The first stage builds the binary.
# The second stage executes the binary.

# Stage 1: Building
# Use the Go Alpine image as building environment
FROM docker.io/library/golang:alpine as builder
# Add user to passwd for stage 2
RUN adduser -HDs /bin/false noroot
# Install ca-certificates for stage and LLVM to strip the binary
RUN apk add ca-certificates llvm
WORKDIR /build
COPY main.go /build/
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o rdns-go main.go
RUN llvm-strip /build/rdns-go

# Stage 2: Running
# Run in minimal environment
FROM scratch
ENV LANG=C.UTF-8
# Copy passwd with our user
COPY --from=builder /etc/passwd /etc/passwd
# Copy the stripped binary
COPY --from=builder --chown=noroot /build/rdns-go /rdns-go
# Execute as non root user
USER noroot
ENTRYPOINT ["/rdns-go"]
