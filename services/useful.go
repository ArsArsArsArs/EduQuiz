package services

import (
	"math/rand"
	"sort"
	"strings"
	"unicode"
	"unicode/utf8"
)

func RemoveEndSymbols(s string, count int) string {
	if count > utf8.RuneCountInString(s) {
		return ""
	}
	return s[:len(s)-count]
}

func CutString(s string, count int) string {
	symbolsAmount := utf8.RuneCountInString(s)
	if count > symbolsAmount {
		return s
	}
	difference := symbolsAmount - count
	return RemoveEndSymbols(s, difference)
}

func CutStringMultipoint(s string, count int) string {
	symbolsAmount := utf8.RuneCountInString(s)
	if count > symbolsAmount {
		return s
	}
	difference := symbolsAmount - count
	return RemoveEndSymbols(s, difference) + "..."
}

type RemovedWord struct {
	Word  string
	Index int
}

func RemoveRandomWords(text string, percent float64) (string, []string) {
	words := strings.Fields(text)
	totalWords := len(words)
	numberToRemove := int(float64(totalWords) * percent / 100)

	removedWordsPre := make([]RemovedWord, 0, numberToRemove)
	removedWords := make([]string, 0, numberToRemove)
	for i := 0; i < numberToRemove; i++ {
		idx := rand.Intn(len(words))
		if words[idx] != "___" {
			removedWordsPre = append(removedWordsPre, RemovedWord{Word: words[idx], Index: idx})
			words[idx] = "___"
		} else {
			i--
		}
	}

	sort.Slice(removedWordsPre, func(i, j int) bool {
		return removedWordsPre[i].Index < removedWordsPre[j].Index
	})
	for _, remWord := range removedWordsPre {
		removedWords = append(removedWords, remWord.Word)
	}

	newText := strings.Join(words, " ")
	return newText, removedWords
}

func CompareStringSlices(a, b []string) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}

func OmitPunctuation(s string) string {
	var result strings.Builder

	for _, r := range s {
		if !unicode.IsPunct(r) {
			result.WriteRune(r)
		}
	}

	return result.String()
}
