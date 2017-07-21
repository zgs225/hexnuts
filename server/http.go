package server

import (
	"encoding/json"
	"html/template"
	"net/http"
	"path/filepath"

	"github.com/zgs225/hexnuts/monitor"
)

type Server struct {
	Configer Configer
	Monitor  *monitor.TCPServer
	Root     http.Dir
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

	s.handleHTML(mux)

	return mux
}

func (s *Server) handleHTML(mux *http.ServeMux) {
	mux.HandleFunc("/ui/", func(w http.ResponseWriter, request *http.Request) {
		t, err := template.ParseFiles(filepath.Join(string(s.Root), "index.html"))
		if err != nil {
			textResponse(w, 500, []byte(err.Error()))
			return
		}
		w.Header().Set("Server", "hexnuts")
		if err := t.Execute(w, nil); err != nil {
			textResponse(w, 500, []byte(err.Error()))
		}
	})

	mux.HandleFunc("/all", func(w http.ResponseWriter, request *http.Request) {
		rv := s.Configer.All()
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json; charset=utf-8")
		w.Header().Set("Server", "hexnuts")
		json.NewEncoder(w).Encode(&rv)
	})

	mux.Handle("/static/", http.FileServer(s.Root))
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
