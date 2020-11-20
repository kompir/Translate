package services

import (
	"errors"
	"fmt"
	"golang.org/x/text/unicode/rangetable"
	"strings"
	"unicode"
)

/** Decoding A Single Word */
func Decode(word string) (string, error) {
	if len(word) > 0 {
		var marks string
		word = strings.ToLower(word)
		/** Remove Marks From Word And Append Them When Is Translated */
		if strings.LastIndexAny(word, "!?,.") == len(word)-1 {
			total := len(word)
			marks = word[total-1 : total]
			word = word[:total-1]
		}
		if strings.ContainsAny(word, "’'") {
			return "", errors.New("please don’t confuse the gophers with apostrophes")
		}
		//** Check If Word Is Already Saved */
		storage := StoredStruct{}
		loaded := storage.ExportWords(word)
		if loaded != "" {
			return fmt.Sprintf("%v", loaded), nil
		}
		words := strings.Split(word, " ")
		/** Translate To Gopher Lang */
		var slices []string
		for _, wrd := range words {
			encoded := translate(wrd)
			slices = append(slices, encoded)
		}
		gopherWord := strings.Join(slices, " ")
		if marks != "" {
			gopherWord += marks
		}
		storage.ImportWords(word, gopherWord)
		return gopherWord, nil
	}
	 return word, errors.New("word is empty")
}

/** Translation Logic For Vowels */
func translate(word string) string {
	vowels := rangetable.New('a', 'e', 'i', 'o', 'u', 'y', 'A', 'E', 'I', 'O', 'U')
	var id int
	for k, v := range word {
		if unicode.Is(vowels, v) {
			if v == 'u' && unicode.Is(vowels, rune(word[k+1])) {
				id++
			}
			if word[:k] == "xr" {
				return "ge" + word[id:] + word[:id]
			}
			if word[:k] == "qu" {
				return word[id:] + word[:id] + "quogo"
			}
			id += k
			break
		}
	}
	if id == 0 {
		return "g" + word
	}
	marks := ""
	if id != 0 {
		marks = "ogo"
	}
	return word[id:] + word[:id] + marks
}

