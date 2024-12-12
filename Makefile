.PHONY: commit push tag release

# Get the current version from git tags
VERSION := $(shell git describe --tags --abbrev=0 2>/dev/null || echo "v0.0.0")
# Increment patch version
NEXT_VERSION := $(shell echo $(VERSION) | awk -F. '{$$NF = $$NF + 1;} 1' | sed 's/ /./g')

# Commit changes with a message
commit:
	@if [ -z "$(m)" ]; then \
		echo "Please provide a commit message using m=<message>"; \
		exit 1; \
	fi
	git add .
	git commit -m "$(m)"

# Create a new tag with incremented version
tag:
	@echo "Current version: $(VERSION)"
	@echo "Next version: $(NEXT_VERSION)"
	git tag $(NEXT_VERSION)
	@echo "Created new tag: $(NEXT_VERSION)"

# Push changes and tags to remote
push:
	git push
	git push --tags

# Combine all steps into one command
release: commit tag push
	@echo "Released version $(NEXT_VERSION)"

# Show help
help:
	@echo "Available commands:"
	@echo "  make commit m=\"your message\"  - Commit changes with a message"
	@echo "  make tag                       - Create a new tag with incremented version"
	@echo "  make push                      - Push changes and tags to remote"
	@echo "  make release m=\"your message\" - Commit, tag, and push in one command" 