package main

import (
	"example/text-perfector/v2/apis"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"

	"github.com/joho/godotenv"
)

// words api
// https://www.wordsapi.com/docs/
// https://www.linguarobot.io/#pricing
// https://dictionaryapi.dev/
//safe text
// https://rapidapi.com/bacloud14/api/safe-text?ref=public_apis&utm_medium=website

//fun translations
// https://api.funtranslations.com/

type Output struct {
	Title   string
	Results any
}

var outputs []Output

func getData(sentence string) map[string]any {
	return map[string]any{
		"Text":    sentence,
		"Words":   strings.Split(strings.ReplaceAll(sentence, ",", " "), " "),
		"Outputs": outputs,
	}
}

var currentSentence = `No man is an island,

Entire of itself,

Every man is a piece of the continent,

A part of the main.

If a clod be washed away by the sea,

It tolls for thee.`

func main() {
	godotenv.Load()
	home := func(w http.ResponseWriter, r *http.Request) {
		templ := template.Must(template.ParseFiles("index.html"))
		templ.Execute(w, getData(currentSentence))
	}
	make_query := func(w http.ResponseWriter, r *http.Request) {
		wordsApi := apis.WordApi{}
		word := r.URL.Query()["word"][0]
		fmt.Println("looking for synonym for", word)
		words := wordsApi.ListSynonyms(word)
		newOutput := Output{Title: fmt.Sprintf("Synonyms for '%s'", word), Results: words}
		outputs = append(
			outputs,
			newOutput,
		)
		fmt.Println(outputs)

		tmpl := template.Must(template.ParseFiles("index.html"))
		tmpl.ExecuteTemplate(w, "output-list-element", newOutput)
	}
	http.HandleFunc("/", home)
	http.HandleFunc("/clicked", make_query)

	// http.HandleFunc("/add-todo", addTodo)
	fmt.Println("Server is running on port 4000: http://localhost:4000")
	log.Fatal(http.ListenAndServe(":4000", logRequest(http.DefaultServeMux)))

}
func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}
