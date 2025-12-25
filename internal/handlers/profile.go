package handlers

import "net/http"

func GetProfile() http.HandlerFunc {	
	return func (w http.ResponseWriter, r *http.Request)  {
		if r.Method != http.MethodGet  {
			http.Error(w, "Bad method", http.StatusBadRequest)
			return 
		}

		w.Header().Set("Content-type", "application/json")
		w.WriteHeader(http.StatusOK)
	}
}