package main

import (
	"fmt"
	"os"
	"strconv"
	"sync"
	"github.com/kompir/Translate/server"
	"github.com/kompir/Translate/services"
)

/** Initiate Server And Storage */
func main() {

	if os.Args[1] == "-port" && len(os.Args) <= 3 && os.Args[2] != "" {
		val, err := strconv.Atoi(os.Args[2])
		if err != nil {
			fmt.Printf("Entered value %s is not a valid port number, please enter correct port.\n", val)
		}
		/** Init Saved Words Storage */
		if services.Saved == nil {
			init := &services.StoredStruct{StoredWords:  nil, StoredSentence: nil, Mutex: sync.RWMutex{}}
			services.Saved = init
		}
		/** Start HTTP Server */
		server.StartServer(":" + strconv.Itoa(val))
	} else {
		fmt.Println("Enter one argument for Port.")
	}
}
