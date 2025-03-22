package main

import (
	"encoding/json"
	"example/text-perfector/v2/apis"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
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

// Postman collection: https://web.postman.co/workspace/Text-App~25f9c420-1881-48ff-9cbb-1f5112f7f870/request/27741934-4e28d433-ed24-443c-b32c-34c54c58284d?tab=headers

type Output struct {
	Title       string
	Synonyms    []string
	Definitions []string
	CleanedText *string
	Language    *string
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

// Example sentence
var currentSentence = `No man is an island,

Entire of itself

Every man is a piece of the continent

A part of the ass.

If a clod be washed away by the sea,

It tolls for thee.`

func main() {
	godotenv.Load()
	home := func(w http.ResponseWriter, r *http.Request) {
		templ := template.Must(template.ParseFiles("index.html"))
		templ.Execute(w, getData(currentSentence))
	}
	makeWordQuery := func(w http.ResponseWriter, r *http.Request) {
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
	makeTextQuery := func(w http.ResponseWriter, r *http.Request) {
		models := []apis.SafeTextModel{}
		if r.PostFormValue("remove-badwords") == "on" {
			models = append(models, apis.RemoveBadWords)
		}
		if r.PostFormValue("punctuation") == "on" {
			models = append(models, apis.Punctuation)
		}
		if r.PostFormValue("detect-language") == "on" {
			models = append(models, apis.DetectLanguage)
		}
		if r.PostFormValue("strip-tags") == "on" {
			models = append(models, apis.StripTags)
		}
		text := r.PostFormValue("input-text")
		fmt.Println("Cleaning text:", text)
		cleanedText, language := apis.CleanText(text, models)
		fmt.Println("Cleaned text:", cleanedText)
		fmt.Println("Language:", language)
		modelsStr := funk.Map(models, func(model apis.SafeTextModel) string {
			return string(model)
		}).([]string)
		newOutput := Output{
			Title:       "Text reformatted with models: " + strings.Join(modelsStr, ", "),
			CleanedText: &cleanedText,
			Language:    language,
		}
		outputs = append(
			outputs,
			newOutput,
		)

		tmpl := template.Must(template.ParseFiles("index.html"))
		tmpl.ExecuteTemplate(w, "output-list-element", newOutput)
	}

	http.HandleFunc("/", home)
	http.HandleFunc("/word-query", makeWordQuery)
	http.HandleFunc("/text-query", makeTextQuery)

	http.HandleFunc("/api/word-query", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {

			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}
		apiKey := r.Header.Get("x-api-Key")
		if apiKey != os.Getenv("api_key") {
			http.Error(w, "API key is missing", http.StatusUnauthorized)
			return
		}

		word := r.URL.Query()["word"][0]
		synonymsChan := make(chan []string)
		definitionsChan := make(chan []string)
		go func() {
			fmt.Println("Looking for synonym for", word)
			synonymsChan <- apis.ListSynonyms(word)
		}()
		go func() {
			fmt.Println("Looking for definition for", word)
			definitionsChan <- apis.GetDefinition(word)
		}()
		synonyms := <-synonymsChan
		definitions := <-definitionsChan
		type Output struct {
			Synonyms    []string `json:"synonyms"`
			Definitions []string `json:"definitions"`
		}
		newOutput := Output{
			Synonyms:    synonyms[:min(3, len(synonyms))],
			Definitions: definitions[:min(3, len(definitions))],
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(newOutput)
	})

	http.HandleFunc("/api/text-query", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {

			http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
			return
		}
		apiKey := r.Header.Get("x-api-Key")
		if apiKey != os.Getenv("api_key") {
			http.Error(w, "API key is missing", http.StatusUnauthorized)
			return
		}

		var requestBody struct {
			Text           string `json:"input-text"`
			RemoveBadWords bool   `json:"remove-badwords"`
			Punctuation    bool   `json:"punctuation"`
			DetectLanguage bool   `json:"detect-language"`
			StripTags      bool   `json:"strip-tags"`
		}
		if err := json.NewDecoder(r.Body).Decode(&requestBody); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		models := []apis.SafeTextModel{}
		if requestBody.RemoveBadWords {
			models = append(models, apis.RemoveBadWords)
		}
		if requestBody.Punctuation {
			models = append(models, apis.Punctuation)
		}
		if requestBody.DetectLanguage {
			models = append(models, apis.DetectLanguage)
		}
		if requestBody.StripTags {
			models = append(models, apis.StripTags)
		}
		text := requestBody.Text
		cleanedText, language := apis.CleanText(text, models)
		fmt.Println("Cleaned text:", cleanedText)
		fmt.Println("Language:", language)
		type Output struct {
			CleanedText *string `json:"cleaned_text"`
			Language    *string `json:"language"`
		}

		newOutput := Output{
			CleanedText: &cleanedText,
			Language:    language,
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(newOutput)
	})

	fmt.Println("Server is running on port 4000: http://localhost:4000")
	log.Fatal(http.ListenAndServe(":4000", logRequest(http.DefaultServeMux)))

}
func logRequest(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Printf("%s %s %s\n", r.RemoteAddr, r.Method, r.URL)
		handler.ServeHTTP(w, r)
	})
}
