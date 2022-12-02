export SESSION_KEY=${skey}

.DEFAULT_GOAL := run
.PHONY: run
run:
	@go run cmd/day${day}/main.go --task=${task}
