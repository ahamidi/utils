package utils

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/pprof"

	"github.com/gorilla/mux"
)

// JSON response convenience function
func jsonResponse(w http.ResponseWriter, status int, data interface{}, err error) {
	w.Header().Add("Content-Type", "application/json")

	// Error Response - Return early
	if err != nil {
		jErr, _ := json.Marshal(map[string]interface{}{
			"Error": err.Error(),
		})
		w.WriteHeader(500)
		w.Write(jErr)
	}

	// Try to handle data
	jRes, mErr := json.Marshal(data)
	if mErr != nil {
		jErr, _ := json.Marshal(map[string]interface{}{
			"Data Error": err.Error(),
		})
		w.WriteHeader(500)
		w.Write(jErr)
	}
	w.WriteHeader(status)
	w.Write(jRes)
}

// Nests map (for adding envelope)
func envelope(d interface{}, envelope string) map[string]interface{} {
	return map[string]interface{}{
		envelope: d,
	}
}

// Unpacks map (opposite process of envelope)
func unvelope(d []byte, envelope string) ([]byte, error) {
	var raw map[string]interface{}

	// Need to use a custom JSON decoder in order to handle large ID
	dec := json.NewDecoder(bytes.NewReader(d))
	dec.UseNumber()
	err := dec.Decode(&raw)
	if err != nil {
		return nil, err
	}

	return json.Marshal(raw[envelope])
}

func enableProfiling(r *mux.Router) {
	r.HandleFunc("/debug/pprof", pprof.Index)
	r.HandleFunc("/debug/pprof/profile", pprof.Profile)
	r.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	r.Handle("/debug/block", pprof.Handler("block"))
	r.Handle("/debug/goroutine", pprof.Handler("goroutine"))
	r.Handle("/debug/heap", pprof.Handler("heap"))
	r.Handle("/debug/threadcreate", pprof.Handler("threadcreate"))
}
