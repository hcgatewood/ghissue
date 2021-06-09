PHONY: release

release:
	goreleaser release --rm-dist
