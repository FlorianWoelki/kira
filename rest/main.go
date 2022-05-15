package main

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/florianwoelki/kira/sandbox"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

var address = ":9090"

type requestBody struct {
	Language string `json:"language"`
	Content  string `json:"content"`
}

func execute(w http.ResponseWriter, r *http.Request) {
	eb := requestBody{}

	if err := json.NewDecoder(r.Body).Decode(&eb); err != nil {
		log.Fatal("an error occurred while decoding json body of `execute` post request", err)
		http.Error(w, "Error decoding json body", http.StatusBadRequest)
		return
	}

	var lang *sandbox.Language
	for _, l := range sandbox.Languages {
		if eb.Language == l.Name {
			lang = &l
			break
		}
	}

	if lang == nil {
		log.Fatalf("no language found with name %s", eb.Language)
		http.Error(w, "Error trying to find valid sandbox runner", http.StatusBadRequest)
		return
	}

	s, output, err := sandbox.Run(lang, eb.Content, []sandbox.SandboxFile{}, []sandbox.SandboxFile{})
	if err != nil {
		log.Fatalf("error while executing sandbox runner: %s", err)
		http.Error(w, "Error trying to execute sandbox runner", http.StatusInternalServerError)
		return
	}

	response, err := json.Marshal(output)
	if err != nil {
		log.Fatalf("error while marshaling output into json: %s", err)
		http.Error(w, "Error trying to marshal output into json", http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	w.Write(response)

	log.Printf("Successful response sent: %s", string(response))

	go func() {
		log.Printf("Cleaning up sandbox with container id %s\n", s.ContainerID)
		s.Clean()
		log.Printf("Cleaned up sandbox with container id %s\n", s.ContainerID)
	}()
}

func loadOrigins(str string) []string {
	result := strings.Split(str, ",")
	for i := 0; i < len(result); i++ {
		result[i] = strings.TrimSpace(result[i])
	}

	return result
}

func main() {
	err := godotenv.Load("local.env")
	if err != nil {
		log.Fatalf("Error occurred while loading env file: %s", err)
	}

	origins := loadOrigins(os.Getenv("ORIGINS"))

	router := mux.NewRouter().StrictSlash(true)
	router.HandleFunc("/execute", execute).Methods(http.MethodPost)

	headersOk := handlers.AllowedHeaders([]string{"X-Requested-With"})
	originsOk := handlers.AllowedOrigins(origins)
	methodsOk := handlers.AllowedMethods([]string{"GET", "POST", "PUT"})

	server := http.Server{
		Addr:         address,
		Handler:      handlers.CORS(originsOk, headersOk, methodsOk)(router),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		log.Println("Starting server on port", address)

		err := server.ListenAndServe()
		if err != nil {
			log.Fatalf("Error starting server: %s\n", err)
			os.Exit(1)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	signal.Notify(c, syscall.SIGTERM)

	sig := <-c
	log.Println("Got signal:", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	server.Shutdown(ctx)
}
