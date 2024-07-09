package main

import "github.com/christo-andrew/haven/cmd/server"

//	@title			Maybe Finance API
//	@version		1.0
//	@description	Maybe Finance API

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		localhost:8080
//	@BasePath	/api/v1

func main() {
	server.StartApp()
}
