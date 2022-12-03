export SESSION_KEY=${skey}

.DEFAULT_GOAL := go
.PHONY: go
go:
	@go run golang/cmd/day${day}/main.go --day=${day} --task=${task}

.PHONY: rust
rust:
	@cargo run -q -p aoc-22-day-${day} -- --day=${day} --task=${task}
