package api

import (
	"log"
	"time"

	"github.com/orangeseeds/rtdatabase/core"
)


func simStream() <-chan string {
	stream := make(chan string, 2)

	data := []string{
		"a",
		"b",
		"c",
		"d",
		"e",
	}

	go func() {
		for _, d := range data {
			time.Sleep(time.Second * 2)
			stream <- d
		}
		close(stream)
	}()

	return stream
}


func printSQliteVersion(db core.Storage) {
	var result string
	row := db.DB.QueryRow("select sqlite_version()")
	_ = row.Scan(&result)

	log.Println(result)
}
