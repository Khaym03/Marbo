package runtime

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"io"
	"os"
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
