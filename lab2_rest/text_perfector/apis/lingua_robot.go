package apis

import (
	"encoding/json"
	"net/http"
	"os"
)

func GetDefinition(word string) []string {
	client := http.Client{}
	req, err := http.NewRequest(
		"GET",
		"https://lingua-robot.p.rapidapi.com/language/v1/entries/en/"+word,
		nil,
	)
	if err != nil {
		panic(err)
	}
	req.Header.Add("x-rapidapi-key", os.Getenv("rapidapi_key"))
	req.Header.Add("x-rapidapi-host", "lingua-robot.p.rapidapi.com")
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != 200 {
		return []string{}
	}
	var data map[string]any
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		panic(err)
	}

	entries := data["entries"].([]any)
	if len(entries) == 0 {
		return []string{}
	}
	item := entries[0].(map[string]any)
	if item["lexemes"] == nil {
		return []string{}
	}
	definitions := []string{}
	lexemes := item["lexemes"].([]any)
	for _, lexeme := range lexemes {
		senses := lexeme.(map[string]any)["senses"].([]any)
		for _, sense := range senses {
			definitions = append(definitions, sense.(map[string]any)["definition"].(string))
		}
	}
	return definitions
}
