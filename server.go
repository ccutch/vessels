package vessels

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Server struct {
	w   http.ResponseWriter
	r   *http.Request
	err error
}

func (u Universe) handle(w http.ResponseWriter, r *http.Request) (*Server, Inventory, Vessel) {
	var (
		i     Inventory
		v     Vessel
		parts = strings.Split(r.URL.Path[1:], "/")
	)
	if len(parts) > 0 {
		i = u[parts[0]]
	}
	if len(parts) > 1 {
		v = i.Get(parts[1:]...)
	}
	return &Server{w, r, nil}, i, v
}

// get is a getter method for getting an inventory from the
// Universe receiver. This will panic if the universe isn't
// registered.
func (s *Server) serveUniverse(u Universe) {
	s.err = json.NewEncoder(s.w).Encode(&u)
}

// server is the proceedure for serving http requests after
// we have parsed and initialized our data.
func (s *Server) serveInventory(i Inventory) {
	switch s.r.Method {
	case "GET":
		s.err = json.NewEncoder(s.w).Encode(query(i, s.r.URL.Query()))
	case "POST":
		v := i.Create()
		_, s.err = io.Copy(io.MultiWriter(s.w, v), s.r.Body)
	default:
		s.err = fmt.Errorf("Unsuppored action: %s", s.r.Method)
	}
}
func (s *Server) serveVessel(v Vessel) {
	switch s.r.Method {
	case "GET":
		_, s.err = io.Copy(s.w, v)
	case "PUT":
		_, s.err = io.Copy(io.MultiWriter(s.w, v), s.r.Body)
	default:
		s.err = fmt.Errorf("Unsupported action: %s", s.r.Method)
	}
}

func (s *Server) serveError() {
	// We are assuming here that if we return an error
	// it is because we looked ahead and found an error
	// with the input. Any errors at runtime will be
	// paniced and handled with recovery.
	if s.err != nil {
		http.Error(s.w, s.err.Error(), http.StatusBadRequest)
	}
}
