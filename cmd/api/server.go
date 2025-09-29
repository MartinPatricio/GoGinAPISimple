package main

import (
	"firstapi/pkg/db"
	"fmt"

	_ "github.com/lib/pq"
)

func main() {
	fmt.Println("Iniciando server")
	if !db.DBStatus() {
		fmt.Println("Ocurrio un error en el ping")
	}
}
