# Build a statically linked binary (Linux):
CGO_ENABLED=0 GOOS=linux go build -ldflags '-s -w -linkmode external -extldflags "-static"' -o rl-bb-static
