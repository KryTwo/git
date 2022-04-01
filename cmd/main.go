package main

import "main/pkg/server"

//открыть подключение к бд хост, название, порт...
//зависимости??????
//http сервер Run()

func main() {

	srv := new(server.Server)
	srv.Run()
}
