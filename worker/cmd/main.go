package main

import (
	"log/slog"
	"os"

	"github.com/nakshatraraghav/transcodex/worker/internal/application"
)

func main() {
	app, err := application.NewApp()
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	err = app.Run()
	if err != nil {
		slog.Error(err.Error())
	}
}
