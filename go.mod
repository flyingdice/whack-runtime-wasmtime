module github.com/flyingdice/whack-runtime-wasmtime

require (
	github.com/bytecodealliance/wasmtime-go v0.33.1
	github.com/flyingdice/whack-sdk v0.0.0-20211208150246-a276f93d7b2a
	github.com/pkg/errors v0.9.1
)

require github.com/google/uuid v1.3.0 // indirect

replace github.com/flyingdice/whack-sdk v0.0.0-20211208150246-a276f93d7b2a => /Users/hawker/src/github.com/flyingdice/whack-sdk

go 1.17
