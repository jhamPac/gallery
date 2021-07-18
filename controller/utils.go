package controller

import (
	"net/http"

	"github.com/gorilla/schema"
)

func parseForm(r *http.Request, dst interface{}) error {
	// Just pointing out issues with implicit state vs pure functions in Haskell
	// Problem with this is the ambiguity. It is not obvious that r has this method or that the result is saved onto a property called PostForm.
	// it is all tribual knowledge, you would have to know that about the request object.
	if err := r.ParseForm(); err != nil {
		panic(err)
	}
	dec := schema.NewDecoder()
	if err := dec.Decode(dst, r.PostForm); err != nil {
		panic(err)
	}
	return nil
}
