package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/chai2010/webp"
)

const VERSION = "0.0.1"

type Config struct {
	Quality       int    `json:"quality"`
	CacheFolder   string `json:"cacheFolder"`
	CacheClearKey string `json:"cacheClearKey"`
	Port          string `json:"port"`
	Debug         bool   `json:"debug"`
}

var (
	configFile = flag.String("config", "config.json", "Path to configuration file")
	config     Config
	debugLog   *log.Logger
)

func loadConfig() error {
	file, err := os.Open(*configFile)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	return decoder.Decode(&config)
}

func setupLogging() {
	if config.Debug {
		debugLog = log.New(os.Stdout, "DEBUG: ", log.Ldate|log.Ltime|log.Lshortfile)
	}
}

func debugLogf(format string, v ...interface{}) {
	if debugLog != nil {
		debugLog.Printf(format, v...)
	}
}

func imageHandler(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Path
	cachePath := config.CacheFolder + filePath + ".webp"

	debugLogf("Handling request for: %s", filePath)

	if _, err := os.Stat(cachePath); err == nil {
		debugLogf("Serving cached file: %s", cachePath)
		http.ServeFile(w, r, cachePath)
		return
	}

	file, err := os.Open(filePath)
	if err != nil {
		debugLogf("File not found: %s", filePath)
		http.Error(w, "File not found", http.StatusNotFound)
		return
	}
	defer file.Close()

	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		debugLogf("Error reading file: %v", err)
		http.Error(w, "Error reading file", http.StatusInternalServerError)
		return
	}
	contentType := http.DetectContentType(buffer)
	debugLogf("Detected content type: %s", contentType)

	_, err = file.Seek(0, 0)
	if err != nil {
		debugLogf("Error seeking file: %v", err)
		http.Error(w, "Error reading file", http.StatusInternalServerError)
		return
	}

	var img image.Image
	switch {
	case strings.Contains(contentType, "jpeg"):
		img, err = jpeg.Decode(file)
	case strings.Contains(contentType, "png"):
		img, err = png.Decode(file)
	case strings.HasPrefix(contentType, "image/webp"):
		debugLogf("Serving original WebP file: %s", filePath)
		http.ServeFile(w, r, filePath)
		return
	default:
		debugLogf("Unsupported file type: %s", contentType)
		http.Error(w, "Unsupported file type", http.StatusBadRequest)
		return
	}

	if err != nil {
		debugLogf("Error processing image: %v", err)
		http.Error(w, "Error processing image", http.StatusInternalServerError)
		return
	}

	if err := os.MkdirAll(filepath.Dir(cachePath), 0755); err != nil {
		debugLogf("Error creating cache directory: %v", err)
		http.Error(w, "Error creating cache directory", http.StatusInternalServerError)
		return
	}

	cacheFile, err := os.Create(cachePath)
	if err != nil {
		debugLogf("Error creating cache file: %v", err)
		http.Error(w, "Error creating cache file", http.StatusInternalServerError)
		return
	}
	defer cacheFile.Close()

	w.Header().Set("Content-Type", "image/webp")
	err = webp.Encode(cacheFile, img, &webp.Options{Quality: float32(config.Quality)})
	if err != nil {
		debugLogf("Error encoding image: %v", err)
		http.Error(w, "Error encoding image", http.StatusInternalServerError)
		return
	}

	debugLogf("Successfully converted and cached image: %s", cachePath)
	http.ServeFile(w, r, cachePath)
}

func cacheClearHandler(w http.ResponseWriter, r *http.Request) {
	debugLogf("Cache clear request received")

	if r.URL.Query().Get("key") != config.CacheClearKey {
		debugLogf("Invalid cache clear key provided")
		http.Error(w, "Invalid cache clear key", http.StatusForbidden)
		return
	}

	if err := os.RemoveAll(config.CacheFolder); err != nil {
		debugLogf("Error clearing cache: %v", err)
		http.Error(w, "Error clearing cache", http.StatusInternalServerError)
		return
	}

	debugLogf("Cache cleared successfully")
	if _, err := w.Write([]byte("Cache cleared successfully")); err != nil {
		debugLogf("Error writing response: %v", err)
		http.Error(w, "Error writing response", http.StatusInternalServerError)
		return
	}
}

func main() {
	fmt.Println("Nimage v" + VERSION)
	flag.Parse()

	if err := loadConfig(); err != nil {
		log.Fatal("Error loading configuration:", err)
	}

	setupLogging()

	debugLogf("Configuration loaded: %+v", config)

	if err := os.MkdirAll(config.CacheFolder, 0755); err != nil {
		log.Fatal("Error creating cache folder:", err)
	}

	http.HandleFunc("/", imageHandler)
	http.HandleFunc("/clearcache", cacheClearHandler)

	log.Printf("Starting server on :%s", config.Port)
	if err := http.ListenAndServe(":"+config.Port, nil); err != nil {
		log.Fatal("ListenAndServe:", err)
	}
}
