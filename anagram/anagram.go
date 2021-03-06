package anagram

import (
	"fmt"
	"math/rand"
	"strings"
	"sort"
)

var Words []string

func MakeAnagram() string {
	Word := Words[rand.Intn(len(Words))]
	anagram := Word
	for anagram == Word {
		anagram = anagrams(Word)
	}
	fmt.Println(Word)	
	return anagram
}

func anagrams(w string) string {
    anagram := ""
	word := strings.Split(w, "")
	for len(word) > 0 {
		i := rand.Intn(len(word))
        	el := word[i]
        	word = append(word[:i], word[i+1:]...)	
		anagram += el
	}
    return anagram
}

func in_array(val string, array []string) (ok bool, i int) {
	for i = range array {	
		if ok = array[i] == val; ok {
			return
		}
	}
	return
}

func Check(anagram string, answer string) bool {
    ok, _ := in_array(answer, Words) 
	return ok && isAnagram(anagram, answer)
}

func isAnagram(original string, test string) bool {
	strOut := sortString(strings.ToLower(strings.Join(strings.Fields(original), "")))
	testOut := sortString(strings.ToLower(strings.Join(strings.Fields(test), "")))
	return strOut == testOut
}

func sortString(w string) string {
	s := strings.Split(w, "")
	sort.Strings(s)
	return strings.Join(s, "")
}