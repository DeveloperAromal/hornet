package requestguard

import (
	"errors"
	"net/http"
)

func (l Limits) RequestGuard(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := l.ValidateURL(r); err != nil {
			http.Error(w, err.Error(), http.StatusRequestURITooLong)
			return
		}

		if err := l.ValidateBody(r); err != nil {
			if errors.Is(err, ErrBodyTooLarge) {
				http.Error(w, err.Error(), http.StatusRequestEntityTooLarge)
				return
			}

			http.Error(w, "failed to read request body", http.StatusBadRequest)
			return
		}

		next.ServeHTTP(w, r)
	})
}
