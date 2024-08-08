package server

import (
	"net/http"

	"github.com/nakshatraraghav/transcodex/backend/internal/routes"
)

func (s *Server) routes() {

	s.router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	routes.UserRouter(s.router, s.db)
	routes.SessionRouter(s.router, s.db)
	routes.ApiKeyRouter(s.router, s.db)
}
