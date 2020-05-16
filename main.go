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
	session "github.com/icza/session"
)

type page struct {
	Word string	
	IsSucess bool
	IsError bool
}

func indexHandler(w http.ResponseWriter, r* http.Request) {	
	if r.URL.Path != "/" {
		return
	}

	t, err := template.ParseFiles("templates/index.html")
	if err != nil {
		fmt.Fprintf(w, err.Error())
	}

	isSucess := false
	isError := false
	
	sess := getSession(w, r)

	switch r.Method {
	case "GET":
		sess.SetAttr("Word", strings.ToUpper(anagram.MakeAnagram()))
	case "POST":
		isSucess = anagram.Check(sess.Attr("Word").(string), r.PostFormValue("answer"))
		isError = !isSucess	
	default:
		fmt.Fprintf(w, "default when switch method")
	}	

	template_error := t.Execute(w, &page{ Word: sess.Attr("Word").(string), IsSucess: isSucess, IsError: isError })
	
	if template_error != nil {
		fmt.Fprintf(w, template_error.Error())
	}
}

func main() {
	// For testing purposes, we want cookies to be sent over HTTP too (not just HTTPS):
	session.Global.Close()
	session.Global = session.NewCookieManagerOptions(session.NewInMemStore(), &session.CookieMngrOptions{AllowHTTP: true})

	content, err := io.ReadFile("dict.txt")
	
	if err != nil {
		panic(err)
	}
	
	rand.Seed(time.Now().UTC().UnixNano())
	anagram.Words = strings.Split(string(content), "\r\n")
	
	http.Handle("/css/", http.StripPrefix("/css/", http.FileServer(http.Dir("css"))))
	http.HandleFunc("/", indexHandler)

	http.ListenAndServe(getPort(), nil)
}

// Get the Port from the environment so we can run on Heroku
func getPort() string {
	port := os.Getenv("PORT")
	fmt.Println(port)

	// Set a default port if there is nothing in the environment
	if port == "" {
		port = "8080"
		fmt.Println("INFO: No PORT environment variable detected, defaulting to " + port)
	}
	return ":" + port
}

func getSession(w http.ResponseWriter, r* http.Request) session.Session {
	sess := session.Get(r)

	if sess == nil {
		sess = session.NewSessionOptions(&session.SessOptions{
			Attrs: map[string]interface{}{"Word": ""},
		})
		session.Add(sess, w)
	} 
	return sess
}