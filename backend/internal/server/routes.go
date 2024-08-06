package server

import "github.com/nakshatraraghav/transcodex/backend/internal/routes"

func (s *Server) routes() {
	routes.UserRouter(s.router)
}
