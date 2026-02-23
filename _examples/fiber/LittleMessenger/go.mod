module LittleMessenger

go 1.25.5



require (
	github.com/XeshSufferer/husocket/core v0.0.0 // indirect
	github.com/XeshSufferer/husocket/adapters/fiber v0.0.0
)

replace (
	github.com/XeshSufferer/husocket/core => ../../../core
	github.com/XeshSufferer/husocket/adapters/fiber => ../../../adapters/fiber
)
