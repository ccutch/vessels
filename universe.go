package vessels

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// Universe is our global namespace of inventories that
// we can represent as a map of string to Inventory.
type Universe map[string]Inventory

// Reference to an inventory and vessel in a universe
type reference struct {
	i Inventory
	v Vessel
}

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
func (u Universe) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// The root index is going to print out the universe
	// as json. This is going to be helpful for debugging
	// but we must ensure this does not expose sensitive
	// information.
	if r.URL.Path == "/" {
		json.NewEncoder(w).Encode(&u)
		return
	}

	ref := u.lookup(r.URL.Path)
	if err := u.serve(w, r, ref); err != nil {
		// We are assuming here that if we return an error
		// it is because we looked ahead and found an error
		// with the input. Any errors at runtime will be
		// paniced and handled with recovery.
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

// lookup is a convenience method for turning a path
// into an inventory and a pointer to a vessel if provided.
func (u Universe) lookup(p string) (ref reference) {
	var ps = strings.Split(p, "/")
	ref.i = u.get(ps[1])
	if len(ps) > 2 {
		ref.v = ref.i.Get(ps[2:]...)
	}
	return
}

// get is a getter method for getting an inventory from the
// Universe receiver. This will panic if the universe isn't
// registered.
func (u Universe) get(n string) (i Inventory) {
	var ok bool
	if i, ok = u[n]; !ok {
		panic(fmt.Sprintf("Inventory not found: %v", n))
	}
	return
}

// server is the proceedure for serving http requests after
// we have parsed and initialized our data.
func (u Universe) serve(w io.Writer, r *http.Request, ref reference) (err error) {
	if ref.v == nil {
		err = json.NewEncoder(w).
			Encode(u.query(ref.i, r.URL.Query()))
		return
	}
	switch r.Method {
	case "GET":
		io.Copy(w, ref.v)
	case "PUT":
		io.Copy(io.MultiWriter(w, ref.v), r.Body)
	default:
		err = fmt.Errorf("Unsupported action: %s", r.Method)
	}
	return
}
