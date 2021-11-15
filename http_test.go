package cache

import (
	"fmt"
	"log"
	"net/http"
	"testing"
)

func TestHTTP(t *testing.T) {
	NewGroup("scores", 2<<10, GetterFunc(
		func(key string) ([]byte, error) {
			log.Println("[SlowDB] search key", key)
			if v, ok := db[key]; ok {
				return []byte(v), nil
			}
			return nil, fmt.Errorf("%s not exist", key)
		}))

	addr := "localhost:9999"
	peers := NewHTTPPool(addr)
	log.Println("minicache is running at", addr)
	log.Fatal(http.ListenAndServe(addr, peers))
}

// curl http://localhost:9999/_minicache/scores/Tom - 630
// curl http://localhost:9999/_minicache/scores/kkk - kkk not exist
