package server

import (
	"net/http"

	"git.youplus.cc/tiny/hexnuts/monitor"
)

type Server struct {
	Configer Configer
	Monitor  *monitor.TCPServer
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

		s.Notify(&monitor.Event{T: monitor.Events_ADD, K: k, V: v})

		textResponse(w, http.StatusOK, []byte("ok"))
	})

	mux.HandleFunc("/update", func(w http.ResponseWriter, request *http.Request) {
		if request.Method != http.MethodPost {
			http.NotFound(w, request)
			return
		}

		request.ParseForm()
		k := request.PostFormValue("key")
		v := request.PostFormValue("value")
		err := s.Configer.Update(k, v)

		if err != nil {
			textResponse(w, http.StatusBadRequest, []byte(err.Error()))
			return
		}

		s.Notify(&monitor.Event{T: monitor.Events_Update, K: k, V: v})

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

		s.Notify(&monitor.Event{T: monitor.Events_DEL, K: k})

		textResponse(w, http.StatusOK, []byte("ok"))
	})

	return mux
}

func (s *Server) Notify(e *monitor.Event) {
	if s.Monitor != nil {
		s.Monitor.Notify(e)
	}
}

func textResponse(w http.ResponseWriter, code int, data []byte) {
	w.WriteHeader(code)
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.Header().Set("Server", "hexnuts")
	w.Write(data)
}
