package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

func getScore(question string, answer string, part int, githubToken string) string {
	url := "https://models.github.ai/inference/chat/completions"
	var sentenceRange string
	switch part {
	case 1:
		sentenceRange = "2–4 short"
	case 2:
		sentenceRange = "8–12 extended"
	case 3:
		sentenceRange = "4–8 analytical"
	}

	// Define the request body structure
	payload := map[string]interface{}{
		"messages": []map[string]string{
			{
				"role": "system",
				"content": fmt.Sprintf(`You are a certified IELTS Speaking examiner. Evaluate this Part %d IELTS Speaking answer using the IELTS Speaking criteria.
				- This is **Part %d**. It must have %s sentences.
				- Part 1: 2–4 short sentences
				- Part 2: 8–12 extended sentences
				- Part 3: 4–8 analytical sentences
			  - If the number of sentences is too low for this part, reduce the **Fluency and Coherence** score by at least 1–2 bands, even if the language is good.
			  
Evaluate based on the 4 IELTS criteria:
1. Fluency and Coherence
2. Lexical Resource
3. Grammatical Range and Accuracy
4. Pronunciation

For each criterion:
- Give a band score (0–9)
- Give a **short reason** (1 sentence only)

Output Format (exactly):
Fluency and Coherence: [score] – [short reason]
Lexical Resource: [score] – [short reason]
Grammatical Range and Accuracy: [score] – [short reason]
Pronunciation: [score] – [short reason]
Overall Band Score: [score]

Keep the output brief and structured. No extra commentary.`, part, part, sentenceRange),
			},
			{
				"role": "user",
				"content": fmt.Sprintf(`
				Question: %s
				Candidate’s Answer:
				%s
				`, question, answer),
			},
		},
		"model":       "openai/gpt-4o-mini",
		"temperature": 0.3,
		"max_tokens":  500,
		"top_p":       0.8,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		panic(err)
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+githubToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	// Read and print response
	bodyBytes, _ := io.ReadAll(resp.Body)
	var result map[string]interface{}
	json.Unmarshal(bodyBytes, &result)

	// Log the raw response (useful for debugging)
	fmt.Println("Raw response:", string(bodyBytes))

	choices, ok := result["choices"].([]interface{})
	if !ok || len(choices) == 0 {
		log.Println("API Error or Unexpected Response Format:", string(bodyBytes))
		return "Error: invalid response from scoring engine"
	}

	message, ok := choices[0].(map[string]interface{})["message"].(map[string]interface{})
	if !ok {
		log.Println("Error reading message field:", string(bodyBytes))
		return "Error: invalid message format"
	}

	content, ok := message["content"].(string)
	if !ok {
		log.Println("Content missing in message:", string(bodyBytes))
		return "Error: content missing"
	}

	return content

}
