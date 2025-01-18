
# Set GOOS and GOARCH to compile for Linux
$env:GOOS = "linux"
$env:GOARCH = "arm64"


# Build the Go binary for Linux and output to a specific file
go build -o fit-byte main.go
