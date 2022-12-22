THIS := $(lastword $(MAKEFILE_LIST))

.PHONY: test

build:
	scripts/build.sh

run:
	bin/ethix

test:
	scripts/test.sh

go:
	@$(MAKE) -f $(THIS) build
	@$(MAKE) -f $(THIS) run
