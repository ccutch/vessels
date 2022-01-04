package vessels

import (
	"net/http"
)

// Universe is our global namespace of inventories that
// we can represent as a map of string to Inventory.
type Universe map[string]Inventory

// NewUniverse makes a new make and returns a universe
// for the newly created map.
func NewUniverse() Universe {
	return Universe(make(map[string]Inventory))
}

// Inventory registers an inventory in the universe.
func (u Universe) Inventory(n string, i Inventory) {
	if _, ok := u[n]; ok {
		panic("Duplicate inventory name " + n)
	}

	// If we see that this universe is an Initializer we
	// will do the initialization now, passing a reference
	// to the Universe that it is being registered to.
	u.initialize(i)
	u[n] = i
}

// ServeHTTP allows us to use a Universe as an http.Handler.
// We will serve information about the universe, inventories
// and vessels via this interface to apps using this service.
func (uni Universe) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s, i, v := uni.handle(w, r)
	if i == nil {
		s.serveUniverse(uni)
	} else if v == nil {
		s.serveInventory(i)
	} else {
		s.serveVessel(v)
	}

	s.serveError()
}
