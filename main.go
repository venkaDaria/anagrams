package main

import (
	"os"
	io "io/ioutil"
	"fmt"
	"time"
    	"html/template"
    	"net/http"
	"strings"
	"math/rand"
	anagram "anagrams/anagram"
)

var Word string
var count int

type page struct {
	Word string	
	IsSucess bool
	IsError bool
}

func indexHandler(w http.ResponseWriter, r* http.Request) {	
	t, err := template.ParseFiles("templates/index.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}
	
	isSucess := false
	isError := false
	
	if count == 0 {
		switch r.Method {
		case "GET":
			Word = strings.ToUpper(anagram.MakeAnagram())
		case "POST":
			isSucess = anagram.Check(Word, r.PostFormValue("answer"))
			isError = !isSucess	
		default:
			fmt.Fprintf(w, "default when switch method")
		}	
	}
	
    	template_error := t.Execute(w, &page{ Word: Word, 
		IsSucess: isSucess, IsError: isError })
	
	if template_error != nil {
		fmt.Fprintf(w, template_error.Error())
	}
	
	count++
	if count == 3 {
		count = 0
	}
}

func main() {
	content, err := io.ReadFile("dict.txt")
	
	if err != nil {
		panic(err)
	}
	
	rand.Seed(time.Now().UTC().UnixNano())
	anagram.Words = strings.Split(string(content), "\n")
	
    	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
    	http.HandleFunc("/", indexHandler)

    	http.ListenAndServe(getPort(), nil)
}

// Get the Port from the environment so we can run on Heroku
func getPort() string {
	var port = os.Getenv("PORT")
	fmt.Println(port)

	// Set a default port if there is nothing in the environment
	if port == "" {
		port = "8080"
		fmt.Println("INFO: No PORT environment variable detected, defaulting to " + port)
	}
	return ":" + port
}