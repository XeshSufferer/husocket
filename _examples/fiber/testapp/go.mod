module testapp

go 1.25.5

require (
	github.com/XeshSufferer/husocket/adapters/fiber v0.0.0
	github.com/XeshSufferer/husocket/core v0.0.0
	github.com/gofiber/fiber/v2 v2.52.11
)


replace (
	github.com/XeshSufferer/husocket/adapters/fiber => ../../../adapters/fiber
	github.com/XeshSufferer/husocket/core => ../../../core
)
