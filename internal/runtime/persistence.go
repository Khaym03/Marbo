package runtime

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"

	"github.com/Khaym03/Marbo/internal/domain"
	"github.com/Khaym03/Marbo/internal/embedder"

	"github.com/Khaym03/Marbo/internal/validator"
)

const CacheFilePath = "cache.json"

func LoadCache(dataFilePath string) (*Cache, error) {
	// 1. Check if cache exists
	if _, err := os.Stat(CacheFilePath); os.IsNotExist(err) {
		return nil, nil // Cache not found
	}

	// 2. Check if cache is stale
	stale, err := isCacheStale(dataFilePath, CacheFilePath)
	if err != nil {
		return nil, err
	}
	if stale {
		return nil, nil // Cache is stale
	}

	// 3. Load cache
	file, err := os.Open(CacheFilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var cache Cache
	if err := json.NewDecoder(file).Decode(&cache); err != nil {
		return nil, err
	}
	return &cache, nil
}

func SaveCache(cache *Cache) error {
	file, err := os.Create(CacheFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	return encoder.Encode(cache)
}

func isCacheStale(dataFilePath, cacheFilePath string) (bool, error) {
	dataHash, err := getFileHash(dataFilePath)
	if err != nil {
		return false, err
	}

	// For simplicity, store hash in a companion file
	hashFilePath := cacheFilePath + ".hash"
	cachedHash, err := os.ReadFile(hashFilePath)
	if os.IsNotExist(err) {
		return true, nil // No hash file, treat as stale
	}
	if err != nil {
		return false, err
	}

	return string(cachedHash) != dataHash, nil
}

func SaveCacheHash(dataFilePath string) error {
	hash, err := getFileHash(dataFilePath)
	if err != nil {
		return err
	}
	return os.WriteFile(CacheFilePath+".hash", []byte(hash), 0644)
}

func getFileHash(filePath string) (string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	hash := sha256.New()
	if _, err := io.Copy(hash, file); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", hash.Sum(nil)), nil
}
func LoadOrBuildCache(
	dataFilePath string,
	emb embedder.Embedder,
) (*Cache, *domain.KnowledgeBase, error) {
	cache, err := LoadCache(dataFilePath)
	if err != nil {
		return nil, nil, err
	}

	data, err := domain.Load(dataFilePath)
	if err != nil {
		return nil, nil, err
	}

	if err := validator.Validate(data); err != nil {
		return nil, nil, err
	}

	if cache != nil {
		log.Println("Cache loaded successfully.")
		return cache, data, nil
	}

	log.Println("Rebuilding cache...")

	b := NewBuilder(emb)
	cache, err = b.Build(data)
	if err != nil {
		return nil, nil, err
	}

	if err := SaveCache(cache); err != nil {
		return nil, nil, err
	}
	if err := SaveCacheHash(dataFilePath); err != nil {
		return nil, nil, err
	}

	log.Println("Cache built and saved.")
	return cache, data, nil
}
