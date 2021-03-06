package server

import (
	"log"
	"net/http"
	"github.com/kompir/Translate/services"
)


func StartServer(port string)  {

	http.HandleFunc("/history", services.History)
	http.HandleFunc("/word", services.TranslateWord)
	http.HandleFunc("/sentence", services.Sentence)
	log.Fatal(http.ListenAndServe(port, nil))
}
