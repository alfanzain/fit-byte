# Set GOOS and GOARCH to compile for Linux
export GOOS="linux"
export GOARCH="arm64"

# Build the Go binary for Linux and output to a specific file
go build -o fit-byte main.go