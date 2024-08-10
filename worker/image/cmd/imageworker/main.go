package imageworker

import (
	"fmt"

	"github.com/nakshatraraghav/transcodex/worker/image/internal/app"
)

func main() {
	app := app.NewImageWorker()

	fmt.Println(app, "created instance of app")
}
