package server

import "net/http"

func (s *Server) healthcheck(w http.ResponseWriter, r *http.Request) {
	s.sendResponse(w, "OK", http.StatusOK)
}
