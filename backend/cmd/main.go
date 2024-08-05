package main

import (
	"log/slog"
	"os"

	"github.com/nakshatraraghav/transcodex/backend/internal/server"
)

func main() {
	server := server.New()

	if err := server.Start(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}
