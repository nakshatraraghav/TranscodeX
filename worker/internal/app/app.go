package app

type MediaProcessor interface {
	ApplyTransformations(map[string]string) ([]byte, error)
}

type Application struct {
}

func NewApp() *Application {
	return &Application{}
}

func (a *Application) Run() {
}
