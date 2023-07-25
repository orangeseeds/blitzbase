package main

import (
	"log"
)

func main() {
	serveCommand := NewServerCommand()
    err  :=serveCommand.Execute()
    if err != nil {
        log.Println(err)
    }
}
