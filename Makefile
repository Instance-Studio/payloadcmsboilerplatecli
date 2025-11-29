APP_NAME := cli
DIST := dist
LDFLAGS := -s -w

PLATFORMS := \
	linux/amd64 \
	linux/arm64

# Default target
all: clean build package

# Clean dist folder
clean:
	rm -rf $(DIST)
	mkdir -p $(DIST)

# Build binaries for all platforms
build:
	@for platform in $(PLATFORMS); do \
		GOOS=$${platform%/*}; \
		GOARCH=$${platform#*/}; \
		echo "Building $$GOOS-$$GOARCH..."; \
		GOOS=$$GOOS GOARCH=$$GOARCH go build -ldflags "$(LDFLAGS)" -o $(DIST)/$(APP_NAME)-$$GOOS-$$GOARCH .; \
	done

# Package binaries into tar.gz
package:
	@for platform in $(PLATFORMS); do \
		GOOS=$${platform%/*}; \
		GOARCH=$${platform#*/}; \
		FILE=$(APP_NAME)-$$GOOS-$$GOARCH; \
		echo "Packaging $$FILE.tar.gz..."; \
		tar -czf $(DIST)/$$FILE.tar.gz -C $(DIST) $$FILE; \
	done

# Remove compiled binaries only
clean-binaries:
	rm -f $(DIST)/$(APP_NAME)-*

.PHONY: all clean build package clean-binaries
