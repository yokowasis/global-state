package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/redis/go-redis/v9"
)

var rdb *redis.Client
var ctx = context.Background()

type GetRequest struct {
	Key string `json:"key"`
}

type SetRequest struct {
	Key   string `json:"key"`
	Value any    `json:"value"`
}

type GetResponse struct {
	Key   string `json:"key"`
	Value any    `json:"value"`
}

type SetResponse struct {
	Success bool   `json:"success"`
	Key     string `json:"key"`
	Value   any    `json:"value"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusOK)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func handleGet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, ErrorResponse{Error: "method not allowed"})
		return
	}

	var req GetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid JSON"})
		return
	}

	val, err := rdb.Get(ctx, req.Key).Result()
	if err == redis.Nil {
		writeJSON(w, http.StatusOK, GetResponse{Key: req.Key, Value: nil})
		return
	}
	if err != nil {
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	var parsed any
	if err := json.Unmarshal([]byte(val), &parsed); err == nil {
		writeJSON(w, http.StatusOK, GetResponse{Key: req.Key, Value: parsed})
		return
	}

	writeJSON(w, http.StatusOK, GetResponse{Key: req.Key, Value: val})
}

func handleSet(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeJSON(w, http.StatusMethodNotAllowed, ErrorResponse{Error: "method not allowed"})
		return
	}

	var req SetRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeJSON(w, http.StatusBadRequest, ErrorResponse{Error: "invalid JSON"})
		return
	}

	var val string
	switch v := req.Value.(type) {
	case string:
		val = v
	default:
		b, err := json.Marshal(v)
		if err != nil {
			writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: "failed to marshal value"})
			return
		}
		val = string(b)
	}

	if err := rdb.Set(ctx, req.Key, val, 0).Err(); err != nil {
		writeJSON(w, http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	writeJSON(w, http.StatusOK, SetResponse{Success: true, Key: req.Key, Value: req.Value})
}

func writeJSON(w http.ResponseWriter, status int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func main() {
	redisAddr := os.Getenv("REDIS_ADDR")
	if redisAddr == "" {
		redisAddr = "localhost:6379"
	}

	rdb = redis.NewClient(&redis.Options{
		Addr: redisAddr,
	})

	if err := rdb.Ping(ctx).Err(); err != nil {
		log.Fatalf("failed to connect to redis: %v", err)
	}
	log.Println("connected to redis")

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		action := r.URL.Query().Get("_")
		switch action {
		case "/get":
			handleGet(w, r)
		case "/set":
			handleSet(w, r)
		default:
			writeJSON(w, http.StatusNotFound, ErrorResponse{Error: fmt.Sprintf("unknown action: %s", action)})
		}
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	handler := corsMiddleware(mux)
	log.Printf("server listening on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}
