package main

func (app *App) NewRPCServer() *RPCServer {
	return NewRPCServer(app.inspectHandler)
}
