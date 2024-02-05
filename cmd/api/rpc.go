package main

import "net/rpc"

func (app *App) RunRPCServer() error {
	if err := rpc.Register(app.rpcServer); err != nil {
		return err
	}
	return app.rpcServer.Listen()
}
