package HTTPServer

import httpserver "github.com/aerosystems/common-service/http_server"

type Config struct {
	httpserver.Config
	mode string
}
