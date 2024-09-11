package server

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/zoturen/seearch/pkg/model"
)

type HttpServer struct {
	indexModel *model.IndexModel
}

func NewHttpServer(im *model.IndexModel) *HttpServer {
	return &HttpServer{
		indexModel: im,
	}
}

func (hs *HttpServer) Serve() {
	http.HandleFunc("/search", enableCORS(hs.searchHandler))

	log.Println("Server starting on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatalf("Server failed to start: %v", err)
	}
}

func (hs *HttpServer) searchHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	query := r.URL.Query().Get("q")
	if query == "" {
		http.Error(w, "Missing query parameter", http.StatusBadRequest)
		return
	}

	searchResult := hs.indexModel.Search(strings.ToLower(query))

	// Set content type
	w.Header().Set("Content-Type", "application/json")

	// Encode and send response
	err := json.NewEncoder(w).Encode(map[string]interface{}{
		"result": searchResult,
	})
	if err != nil {
		http.Error(w, "Error encoding response", http.StatusInternalServerError)
		return
	}
}

func enableCORS(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	}
}
