package main

import (
	"example/text-perfector/v2/apis"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"regexp"
	"slices"
	"strings"

	"github.com/joho/godotenv"
	"github.com/thoas/go-funk"
)

// words api
// https://www.wordsapi.com/docs/
// https://www.linguarobot.io/#pricing
// https://dictionaryapi.dev/
//safe text
// https://rapidapi.com/bacloud14/api/safe-text?ref=public_apis&utm_medium=website

//fun translations
// https://api.funtranslations.com/

// tutorial
// https://www.youtube.com/watch?v=F9H6vYelYyU

type Output struct {
	Title       string
	Synonyms    []string
	Definitions []string
}

var outputs []Output

func getData(sentence string) map[string]any {
	return map[string]any{
		"Text":    sentence,
		"Words":   tokenizeSentence(sentence),
		"Outputs": outputs,
	}
}
func tokenizeSentence(sentence string) []string {
	re := regexp.MustCompile(`[^a-zA-Z]+`)
	words := re.Split(sentence, -1)
	words = funk.Uniq(words).([]string)
	words = funk.FilterString(words, func(word string) bool {
		return len(word) > 3
	})
	words = funk.Map(words, func(word string) string {
		return strings.ToLower(word)
	}).([]string)
	slices.Sort(words)
	return words
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
		word := r.URL.Query()["word"][0]
		fmt.Println("Looking for synonym for", word)
		synonymsList := apis.ListSynonyms(word)
		fmt.Println("Looking for definition for", word)
		definitions := apis.GetDefinition(word)
		newOutput := Output{
			Title:       fmt.Sprintf("Query for word: '%s'", word),
			Synonyms:    synonymsList[:min(3, len(synonymsList))],
			Definitions: definitions[:min(3, len(definitions))],
		}
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
