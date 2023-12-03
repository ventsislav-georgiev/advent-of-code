.DEFAULT_GOAL := go
.PHONY: go
go:
	@go run golang/cmd/${year}/${day}/main.go --task=${task} --input=${input}

.PHONY: rust
rust:
	@cargo run -q -p aoc-${year}-${day} -- --year=${year} --day=${day} --task=${task}

.PHONY: bench
bench:
	@go test -bench=. -benchtime=20s -benchmem -run=^$$ ./golang/cmd/${year}/${day}

.PHONY: test
test:
	@go test -count=1 ./golang/cmd/${year}/${day}

.PHONY: synacor
synacor:
	@go run golang/cmd/synacor/main.go --telecode=${TELECODE} --vault=${VAULT}
