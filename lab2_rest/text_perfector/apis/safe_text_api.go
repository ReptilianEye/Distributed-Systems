package apis

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
)

type SafeTextModel string

const (
	RemoveBadWords SafeTextModel = "BadWords"
	Punctuation    SafeTextModel = "Punctuate"
	DetectLanguage SafeTextModel = "DetectLanguage"
	StripTags      SafeTextModel = "StripTags"
)

func CleanText(text string, models []SafeTextModel) (string, *string) {
	client := http.Client{}
	req, err := http.NewRequest(
		"GET",
		"https://safe-text.p.rapidapi.com/clean_text",
		nil,
	)
	if err != nil {
		panic(err)
	}
	q := req.URL.Query()
	q.Add("text", text)

	modelStrings := make([]string, len(models))
	for i, model := range models {
		modelStrings[i] = string(model)
	}
	q.Add("models", strings.Join(modelStrings, ","))

	req.URL.RawQuery = q.Encode()
	req.Header.Add("x-rapidapi-key", os.Getenv("rapidapi_key"))
	req.Header.Add("x-rapidapi-host", "safe-text.p.rapidapi.com")
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	if resp.StatusCode != 200 {
		return "", nil
	}
	var data map[string]any
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		panic(err)
	}
	result := data["call_response"].(map[string]any)["result"].(map[string]any)
	resultStr := result["clean"].(string)
	if result["stripped"] != nil {
		resultStr = result["stripped"].(string)
	}
	additional := result["additional"].(map[string]any)
	if additional["language"] == nil {
		return resultStr, nil
	}
	language := additional["language"].(string)
	return resultStr, &language

}
