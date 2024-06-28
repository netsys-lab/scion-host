# /!bin/bash

set -e

cd dev/scion

rm -rf build
mkdir build

echo "Building SCION binaries. This may take a while, especially when building the first time."

# Daemon Binary
#echo "Building Daemon..."
#cd daemon/cmd/daemon
#CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o daemon-darwin-amd64
#CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o daemon-darwin-arm64
#CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o daemon
#CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o daemon.exe

#cd ../../../
#cp daemon/cmd/daemon/daemon-darwin-amd64 ../../environment/darwinx64/daemon
#cp daemon/cmd/daemon/daemon-darwin-arm64 ../../environment/darwinarm64/daemon
#cp daemon/cmd/daemon/daemon ../../environment/linuxx64/daemon
#cp daemon/cmd/daemon/daemon.exe ../../environment/windowsx64/daemon.exe

# Dispatcher Binary
#echo "Building Dispatcher..."
#cd dispatcher/cmd/dispatcher
#CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o dispatcher-darwin-amd64
#CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o dispatcher-darwin-arm64
#CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o dispatcher
#CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o dispatcher.exe

#cd ../../../
#cp dispatcher/cmd/dispatcher/dispatcher-darwin-amd64 ../../environment/darwinx64/dispatcher
#cp dispatcher/cmd/dispatcher/dispatcher-darwin-arm64 ../../environment/darwinarm64/dispatcher
#cp dispatcher/cmd/dispatcher/dispatcher ../../environment/linuxx64/dispatcher
#cp dispatcher/cmd/dispatcher/dispatcher.exe ../../environment/windowsx64/dispatcher.exe

# SCION Binary
echo "Building SCION..."
cd scion/cmd/scion 
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o scion-darwin-amd64
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o scion-darwin-arm64
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o scion
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o scion.exe

cd ../../../
cp scion/cmd/scion/scion-darwin-amd64 ../../environment/darwinx64/scion
cp scion/cmd/scion/scion-darwin-arm64 ../../environment/darwinarm64/scion
cp scion/cmd/scion/scion ../../environment/linuxx64/scion
cp scion/cmd/scion/scion.exe ../../environment/windowsx64/scion.exe

cp scion/cmd/scion/scion-darwin-amd64 ../../build/scion-darwin-amd64
cp scion/cmd/scion/scion-darwin-arm64 ../../build//scion-darwin-arm64
cp scion/cmd/scion/scion ../../build/scion
cp scion/cmd/scion/scion.exe ../../build/scion.exe

cd ../../



# SCION Host Binary
echo "Building SCION Host..."
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o build/scion-host-darwin-amd64
CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -o build/scion-host-darwin-arm64
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o build/scion-host
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o build/scion-host.exe

echo "Done"