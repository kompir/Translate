package services

import (
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

var translateWord = []string{
	"apple", "epple",
}

var expectedWords = []string{
	"gapple", "gepple",
}

var wrongWords = []string{
	"apples's", "don’t",
}

var memoryWord = "gapple"

/** Testing Translation Logic Function  */
func TestWord(t *testing.T) {
	if Saved == nil {
		init := &StoredStruct{StoredWords:  nil, StoredSentence: nil, Mutex: sync.RWMutex{}}
		Saved = init
	}
	for k, v := range translateWord {
		result, err := Word(v)
		assert.Nil(t, err)
		assert.Equal(t, expectedWords[k], result["gopher-word"])

	}
}

/** Testing Error Handling */
func TestDecode(t *testing.T) {
	if Saved == nil {
		init := &StoredStruct{StoredWords:  nil, StoredSentence: nil, Mutex: sync.RWMutex{}}
		Saved = init
	}
	for _, v := range wrongWords {
		_, err := Decode(v)
		assert.EqualError(t, err, "please don’t confuse the gophers with apostrophes")
	}
}

/** Test Export From Memory */
func TestExportWords(t *testing.T) {
	storage := StoredStruct{}
	storage.ImportWords("apple", memoryWord)
	exported := storage.ExportWords("apple")
	assert.Equal(t, memoryWord, exported)
}





