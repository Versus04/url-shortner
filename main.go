package main

import (
	"encoding/json"
	"net/http"
	"os"
	"strings"
	"sync"
)

type store struct {
	Counter int64             `json:"counter"`
	URLs    map[string]string `json:"urls"`
}

var mu sync.RWMutex

const chars = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
const base_url = "http://localhost:8080/"

var count int64 = 78992
var mpp = make(map[string]string)

func reverse(b []byte) []byte {
	for i, j := 0, len(b)-1; i < j; i, j = i+1, j-1 {
		b[i], b[j] = b[j], b[i]
	}
	return b
}

func encodeBase62(n int64) string {
	var result []byte
	for n > 0 {
		result = append(result, chars[n%62])
		n /= 62
	}
	return string(reverse(result))
}
func submitURL(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Only POST allowed", http.StatusMethodNotAllowed)
		return
	}
	ct := r.Header.Get("Content-Type")
	if !strings.HasPrefix(ct, "application/json") {
		http.Error(w, "Expected application/json", http.StatusUnsupportedMediaType)
		return
	}

	defer r.Body.Close()

	var req struct {
		URL string `json:"url"`
	}

	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid JSON body", http.StatusBadRequest)
		return
	}
	if req.URL == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var response struct {
		Short_Code string `json:"short_code"`
		Short_URL  string `json:"short_url"`
	}
	mu.Lock()
	response.Short_Code = encodeBase62(int64(count))
	response.Short_URL = base_url + response.Short_Code
	mpp[response.Short_Code] = req.URL
	count++
	saveFile()
	mu.Unlock()
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)

}
func getURL(w http.ResponseWriter, r *http.Request) {
	code := strings.TrimPrefix(r.URL.Path, "/r/")
	if code == "" {
		http.Error(w, "Code is missing", http.StatusBadRequest)
		return
	}
	mu.RLock()
	url, ok := mpp[code]
	mu.RUnlock()
	if !ok {
		http.NotFound(w, r)
		return
	}
	http.Redirect(w, r, url, http.StatusFound)
}
func saveFile() {

	save := store{
		Counter: count,
		URLs:    mpp,
	}
	file, err := json.Marshal(save)
	if err != nil {
		return
	}
	os.WriteFile("store.json", file, 0644)

}
func loadTask() store {
	file, err := os.ReadFile("store.json")
	if err != nil {
		return store{
			Counter: 123456789,
			URLs:    make(map[string]string),
		}
	}
	var datamap store
	err = json.Unmarshal(file, &datamap)
	if err != nil {
		return store{
			Counter: 123456789,
			URLs:    make(map[string]string),
		}
	}
	return datamap
}
func main() {
	store1 := loadTask()
	mpp = store1.URLs
	count = store1.Counter
	http.HandleFunc("/shorten", submitURL)
	http.HandleFunc("/r/", getURL)
	//mpp["huzz"] = "huzz"
	http.ListenAndServe(":8080", nil)

}
