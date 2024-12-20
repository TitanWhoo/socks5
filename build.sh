#!/usr/bin/env sh
set -e
PROJECT_NAME="socks5"
DARWIN_AMD64=./output/${PROJECT_NAME}_darwin_amd64
DARWIN_ARM64=./output/${PROJECT_NAME}_darwin_arm64
LINUX_AMD64=./output/${PROJECT_NAME}_linux_amd64
LINUX_ARM64=./output/${PROJECT_NAME}_linux_arm64
WINDOWS_AMD64=./output/${PROJECT_NAME}_windows_amd64.exe

if [ "$1" = "all" ]; then
    echo "start to build all platforms"
    echo "build darwin amd64 to ${DARWIN_AMD64}"
    CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o ${DARWIN_AMD64} cli/main.go
    echo "build darwin arm64 to ${DARWIN_ARM64}"
    CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o ${DARWIN_ARM64} cli/main.go
    echo "build linux amd64 to ${LINUX_AMD64}"
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o ${LINUX_AMD64} cli/main.go
    echo "build linux arm64 to ${LINUX_ARM64}"
    CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags="-s -w" -o ${LINUX_ARM64} cli/main.go
    echo "build windows amd64 to ${WINDOWS_AMD64}"
    CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o ${WINDOWS_AMD64} cli/main.go
    echo "done"
    exit 0
elif [ "$1" = "deploy" ]; then
    echo "build linux amd64 to ${LINUX_AMD64}"
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o ${LINUX_AMD64} cli/main.go
    upx -9 ${LINUX_AMD64}
elif [ "$1" = "install" ]; then
    go build -ldflags="-s -w" -o /usr/local/${PROJECT_NAME}/${PROJECT_NAME} cli/main.go
elif [ "$1" = "linux" ]; then
    echo "build linux amd64 to ${LINUX_AMD64}"
    CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o ${LINUX_AMD64} cli/main.go
elif [ "$1" = "windows" ]; then
    echo "build windows amd64 to ${WINDOWS_AMD64}"
    CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o ${WINDOWS_AMD64} cli/main.go
elif [ "$1" = "darwin" ]; then
    echo "build darwin amd64 to ${DARWIN_AMD64}"
    CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o ${DARWIN_AMD64} cli/main.go
    echo "build darwin arm64 to ${DARWIN_ARM64}"
    CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags="-s -w" -o ${DARWIN_ARM64} cli/main.go
else
    go build -ldflags="-s -w"  -o ./output/${PROJECT_NAME} cli/main.go
    exit 0
fi