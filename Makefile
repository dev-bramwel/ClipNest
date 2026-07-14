.PHONY: lint

lint:
	@echo "Running standard code formatting checks..."
	@# Capture the list of poorly formatted files
	@files=$$(gofmt -l .); \
	if [ -n "$$files" ]; then \
		echo "Formatting the following files:"; \
		echo "$$files"; \
		gofmt -w .; \
		echo "Files have been auto-formatted."; \
	else \
		echo "Glow green! All code conforms to standard Go formatting specifications."; \
	fi
	@echo "Running basic compiler structural checks..."
	@go vet ./...
	