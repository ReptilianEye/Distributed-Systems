package apis

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
)

func ListSynonyms(word string) []string {
	client := http.Client{}
	req, err := http.NewRequest("GET", "https://wordsapiv1.p.rapidapi.com/words/"+word, nil)
	if err != nil {
		panic(err)
	}
	req.Header.Add("X-Mashape-Key", os.Getenv("rapidapi_key"))
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
	results := data["results"].([]any)
	if len(results) == 0 {
		return []string{}
	}
	result := results[0].(map[string]any)
	if result["synonyms"] == nil {
		return []string{}
	}
	synonyms := result["synonyms"].([]any)
	fmt.Println(synonyms)
	synonymsStr := []string{}
	for _, synonym := range synonyms {
		synonymsStr = append(synonymsStr, synonym.(string))
	}
	return synonymsStr

}
