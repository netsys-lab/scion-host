.PHONY: build dev clean

build:
	@echo "Running build script..."
	@./build.sh

dev:
	@echo "Running development script..."
	@./dev.sh

clean:
	@echo "Cleaning build environment..."
	@./clean.sh