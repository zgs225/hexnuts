package server

import (
	"net/http"
)

type Server struct {
	Configer Configer
}

func (s *Server) MakeHTTPServer() http.Handler {
	mux := http.NewServeMux()

	mux.HandleFunc("/set", func(w http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodPost {
			http.NotFound(w, request)
			return
		}

		request.ParseForm()
		k := request.PostFormValue("key")
		v := request.PostFormValue("value")
		err := s.Configer.Set(k, v)

		if err != nil {
			textResponse(w, http.StatusBadRequest, []byte(err.Error()))
			return
		}

		textResponse(w, http.StatusOK, []byte("ok"))
	})

	mux.HandleFunc("/get", func(w http.ResponseWriter, request *http.Request) {
		request.ParseForm()
		k := request.FormValue("key")
		v, err := s.Configer.Get(k)

		if err != nil {
			textResponse(w, http.StatusNotFound, []byte(err.Error()))
			return
		}

		textResponse(w, http.StatusOK, []byte(v))
	})

	mux.HandleFunc("/del", func(w http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodPost {
			http.NotFound(w, request)
			return
		}

		request.ParseForm()
		k := request.PostFormValue("key")
		err := s.Configer.Del(k)

		if err != nil {
			textResponse(w, http.StatusBadRequest, []byte(err.Error()))
			return
		}

		textResponse(w, http.StatusOK, []byte("ok"))
	})

	return mux
}

func textResponse(w http.ResponseWriter, code int, data []byte) {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Server", "hexnuts")
	w.Write(data)
}
