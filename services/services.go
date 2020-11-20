package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sort"
	"strings"
	"sync"
)

/** Interface For Memory Operations */
type StorageInterface interface {
	ImportWords(englishWord string, gopherWord string)
	ExportWords(word string) string
	ExportAll() map[string]string
}

/** Internal Memory Variable For Saved Data */
var Saved *StoredStruct

/** Structure For Saved Data */
type StoredStruct struct {
	StoredWords map[string]string
	StoredSentence map[string]string
	Mutex sync.RWMutex
}

type englishSentence struct {
	sentence map[string]string
}

type englishWord struct {
	word map[string]string
}

type Combined struct {
	word map[string]string
	sen  map[string]string
}

/** Store Words In Memory */
func (storage *StoredStruct) ImportWords(englishWord string, gopherWord string) {
	storage.Mutex.Lock()
	defer storage.Mutex.Unlock()
	if Saved.StoredWords == nil {
		Saved.StoredWords = make(map[string]string, 0)
	}
	Saved.StoredWords[englishWord] = gopherWord
}

/** Store Sentence In Memory */
func (storage *StoredStruct) ImportSentence(englishSentence string, gopherSentence string) {
	storage.Mutex.Lock()
	defer storage.Mutex.Unlock()
	if Saved.StoredSentence == nil {
		Saved.StoredSentence = make(map[string]string, 0)
	}
	Saved.StoredSentence[englishSentence] = gopherSentence
}

/** Get Words From Memory */
func (storage *StoredStruct) ExportWords(word string) string {
	storage.Mutex.RLock()
	defer storage.Mutex.RUnlock()
	exported := Saved.StoredWords[word]
	return exported
}

/** Get Sentence From Memory */
func (storage *StoredStruct) ExportSentence(sentence string) string {
	storage.Mutex.RLock()
	defer storage.Mutex.RUnlock()
	exported := Saved.StoredSentence[sentence]
	return exported
}

/** Export All Words And Sentences */
func (storage *StoredStruct) ExportAll() map[string]string {
	storage.Mutex.Lock()
	defer storage.Mutex.Unlock()
	/** Export Words Ascending */
	ordersWord := make([]string, 0, len(Saved.StoredWords))
	for word := range Saved.StoredWords {
		ordersWord = append(ordersWord, word)
	}
	sort.Strings(ordersWord)
	ordersSen := make([]string, 0, len(Saved.StoredSentence))
	for sen := range Saved.StoredSentence {
		ordersSen = append(ordersSen, sen)
	}
	sort.Strings(ordersSen)
	exportedSen := make(map[string]string, len(ordersSen))
	for _, sen := range ordersSen {
		exportedSen[sen] = Saved.StoredSentence[sen]
	}
	exportedWord := make(map[string]string, len(ordersWord))
	for _, word := range ordersWord {
		exportedWord[word] = Saved.StoredWords[word]
	}
	for k, v := range exportedSen {
		exportedWord[k] = v
	}
	return exportedWord
}

/** Translate Gopher Word From Request */
func TranslateWord(res http.ResponseWriter, req *http.Request) {
	var englishWord englishWord
	/** Decode Gopher Word From Request*/
	err := json.NewDecoder(req.Body).Decode(&englishWord.word)
	if err != nil {
		http.Error(res, err.Error(), http.StatusBadRequest)
		return
	}
	if len(strings.Fields(englishWord.word["english-word"])) > 1 {
		http.Error(res, "Provide One Word In Format {'english-word':'<a single English word'}", 500)
		return
	}
	resutl, err := Word(englishWord.word["english-word"])
	if err != nil {
		http.Error(res, err.Error(), 500)
	}
	JsonEncode(res, req, resutl)
}

/** Translate A Single Word */
func Word(word string) (map[string]string, error){
	/** Translate Gopher word */
	translated, err := Decode(word)
	if err != nil {
		return nil, err
	}
	/** Format output */
	formatted := map[string]string{"gopher-word": translated}
	return formatted, nil
}

/** Translate Sentence From Request */
func Sentence(res http.ResponseWriter, req *http.Request ) {
	var englishSentence englishSentence
	var gopherSentence []string
	/** Decode Sentence Word From Request*/
	err := json.NewDecoder(req.Body).Decode(&englishSentence.sentence)
	if err != nil {log.Fatal(err)}
	storage := StoredStruct{}
	sen := englishSentence.sentence["english-sentence"]
	loaded := storage.ExportSentence(sen)
	if loaded != "" {
		JsonEncode(res, req,  map[string]string{"gopher-sentence": loaded})
		return
	}
	/** Split Sentence In Words */
	splitSentence := strings.Split(englishSentence.sentence["english-sentence"], " ")
	if !strings.Contains(".?!", string(englishSentence.sentence["english-sentence"][len(englishSentence.sentence["english-sentence"])-1])) {
		fmt.Fprintf(res, "Please provide sentence with dot, question or exclamation mark at the end.")
	}
	//** Translate Each Gopher word */
	for _, sen := range splitSentence {
		trans, _ := Decode(sen)
		gopherSentence = append(gopherSentence, trans)
	}
	gopherSentenceJoined := strings.Join(gopherSentence, " ")
	/** Format output */
	formatted := map[string]string{"gopher-sentence": gopherSentenceJoined}
	storage.ImportSentence(sen, gopherSentenceJoined)
	JsonEncode(res, req, formatted)
}

/** Get All Translated Data From The Start Of Server */
func History(res http.ResponseWriter, req *http.Request) {
	var words StoredStruct
	saved := words.ExportAll()
	JsonEncode(res, req, saved)
}

/** Json Encode Data In Response */
func JsonEncode(res http.ResponseWriter, req *http.Request, v interface{})  {
	buf := &bytes.Buffer{}
	enc := json.NewEncoder(buf)
	enc.SetEscapeHTML(true)
	if err := enc.Encode(v); err != nil {
		http.Error(res, err.Error(), http.StatusInternalServerError)
		return
	}
	res.Header().Set("Content-Type", "application/json; charset=utf-8")
	if status, ok := req.Context().Value("Status").(int); ok {
		res.WriteHeader(status)
	}
	res.Write(buf.Bytes())
}