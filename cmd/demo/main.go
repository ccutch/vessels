package main

import (
	"net/http"
	"vessels"
	"vessels/memory"
)

func main() {
	u := vessels.NewUniverse()
	u.Inventory("memory", memory.Empty())
	http.ListenAndServe(":8080", u)
}
