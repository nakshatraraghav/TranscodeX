package main

import (
	"log/slog"

	"github.com/nakshatraraghav/transcodex/worker/config"
	"github.com/nakshatraraghav/transcodex/worker/internal/processors/image"
)

// import (
// 	"log/slog"
// 	"os"

// 	"github.com/nakshatraraghav/transcodex/worker/internal/application"
// )

// func main() {
// 	app, err := application.NewApp()
// 	if err != nil {
// 		slog.Error(err.Error())
// 		os.Exit(1)
// 	}

// 	err = app.Run()
// 	if err != nil {
// 		slog.Error(err.Error())
// 	}
// }

func main() {
	config.LoadEnv()

	ip := image.NewImageProcessor()

	ip.LoadData()
	err := ip.GenerateThumbnail("200")
	if err != nil {
		slog.Error(err.Error())
	}
	ip.SaveChanges()

}

// IMAGE => libvips
