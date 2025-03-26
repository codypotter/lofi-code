package main

import (
	"log"
)

func main() {
	// Create the search handler
	// handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 	// query := r.URL.Query().Get("q")
	// 	// results, err := search.SearchPosts(context.Background(), query)
	// 	// if err != nil {
	// 	// 	http.Error(w, "Failed to search posts", http.StatusInternalServerError)
	// 	// 	return
	// 	// }

	// 	// w.Header().Set("Content-Type", "application/json")
	// 	// json.NewEncoder(w).Encode(results)
	// })

	// Run the Lambda function
	log.Println("Starting Lambda...")
}
