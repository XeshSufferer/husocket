module testapp

go 1.25.5

require (
	github.com/XeshSufferer/husocket/adapters/stdhttp v0.0.0
	github.com/XeshSufferer/husocket/core v0.0.0
)

replace (
	github.com/XeshSufferer/husocket/adapters/stdhttp => ../../../adapters/stdhttp
	github.com/XeshSufferer/husocket/core => ../../../core
)
