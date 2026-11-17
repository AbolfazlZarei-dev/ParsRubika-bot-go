package forcejoin

import "net/http"

type failing struct{}

func (failing) ServeHTTP(w http.ResponseWriter, _ *http.Request) {
	http.Error(w, "down", http.StatusInternalServerError)
}
