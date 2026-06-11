GO ?= go

run-all:
	@$(MAKE) -j3 run-command run-query run-worker

run-command:
	$(GO) -C packages/product-command run .

run-query:
	$(GO) -C packages/product-query run .

run-worker:
	$(GO) -C packages/product-worker run .
