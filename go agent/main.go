package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
)

type answer struct {
	Question string `json:"question"`
	Answer   string `json:"answer"`
	Part     int    `json:"part"`
}

type app struct {
	addr  string
	token string
}

func (a *app) Score(w http.ResponseWriter, r *http.Request) {
	log.Println("Received /score request")

	var reqData answer
	err := json.NewDecoder(r.Body).Decode(&reqData)
	if err != nil {
		log.Println("Error decoding JSON:", err)
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}
	log.Printf("Request Data: Question=%q, Part=%d\n", reqData.Question, reqData.Part)

	score := getScore(reqData.Question, reqData.Answer, reqData.Part, a.token) // <-- assumes you add token
	log.Println("Score received from model:\n" + score)

	response := map[string]string{
		"score": score,
	}
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(response)
	if err != nil {
		log.Println("JSON encode error:", err)
	}
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, assuming production environment")
	}
	app := &app{addr: ":8080", token: os.Getenv("AI_API")}

	mux := http.NewServeMux()
	mux.HandleFunc("POST /score", app.Score)

	log.Printf("Server starting on %s\n", app.addr)
	server := &http.Server{
		Addr:    app.addr,
		Handler: mux,
	}
	err = server.ListenAndServe()
	if err != nil {
		log.Fatal("Server failed:", err)
	}
}
