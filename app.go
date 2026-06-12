package main

import (
	"context"
	"fmt"
	"github.com/Khaym03/Marbo/internal/runtime"
)

// App struct
type App struct {
	ctx context.Context
	r   *runtime.Runtime
}

// NewApp creates a new App application struct
func NewApp() *App {
	return &App{}
}

func (a *App) SetRuntime(r *runtime.Runtime) {
	a.r = r
}

func (a *App) SendMessage(text string) (runtime.RuntimeResult, error) {
	if a.r == nil {
		return runtime.RuntimeResult{}, fmt.Errorf("runtime not initialized")
	}
	return a.r.Handle(text)
}

// startup is called when the app starts. The context is saved
// so we can call the runtime methods
func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
}

// Greet returns a greeting for the given name
func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}
