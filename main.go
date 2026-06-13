package main

import (
	"embed"
	"log"

	"github.com/Khaym03/Marbo/internal/embedder"
	"github.com/Khaym03/Marbo/internal/runtime"
	"github.com/wailsapp/wails/v2"
	"github.com/wailsapp/wails/v2/pkg/options"
	"github.com/wailsapp/wails/v2/pkg/options/assetserver"
	ort "github.com/yalue/onnxruntime_go"
)

const (
	dataFilepath = "data.json"
)

//go:embed all:frontend/dist
var assets embed.FS

func main() {
	// Initialize Runtime
	ort.SetSharedLibraryPath("third-party/onnxruntime.dll")
	err := ort.InitializeEnvironment()
	if err != nil {
		log.Fatal(err)
	}
	defer ort.DestroyEnvironment()

	emb, err := embedder.New("modelo_e5_onnx/model.onnx", "modelo_e5_onnx/tokenizer.json")
	if err != nil {
		log.Fatal(err)
	}
	defer emb.Close()

	cache, data, err := runtime.LoadOrBuildCache(dataFilepath, emb)
	if err != nil {
		log.Fatal(err)
	}

	app := NewApp()
	app.SetRuntime(runtime.NewRuntime(emb, cache, data.Settings))

	// Create application with options
	err = wails.Run(&options.App{
		Title:  "Marbo",
		Width:  1024,
		Height: 768,
		AssetServer: &assetserver.Options{
			Assets: assets,
		},
		BackgroundColour: &options.RGBA{R: 255, G: 255, B: 255, A: 1},
		OnStartup:        app.startup,
		Bind: []any{
			app,
		},
	})

	if err != nil {
		println("Error:", err.Error())
	}
}
