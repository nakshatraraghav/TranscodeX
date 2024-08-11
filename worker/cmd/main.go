package main

import (
	"fmt"
	"log/slog"
	"os"

	"github.com/nakshatraraghav/transcodex/worker/config"
)

func main() {
	if err := config.LoadEnv(); err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}

	fmt.Println(config.GetEnv())
}
