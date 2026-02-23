module testapp

go 1.25.5

require (
	github.com/XeshSufferer/husocket/adapters/fasthttp v0.0.0
	github.com/XeshSufferer/husocket/core v0.0.0
	github.com/valyala/fasthttp v1.69.0
)

require (
	github.com/andybalholm/brotli v1.2.0 // indirect
	github.com/fasthttp/websocket v1.5.12 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/klauspost/compress v1.18.3 // indirect
	github.com/savsgio/gotils v0.0.0-20240704082632-aef3928b8a38 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	golang.org/x/net v0.48.0 // indirect
)

replace (
	github.com/XeshSufferer/husocket/adapters/fasthttp => ../../../adapters/fasthttp
	github.com/XeshSufferer/husocket/core => ../../../core
)
