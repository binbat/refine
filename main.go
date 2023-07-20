package main

import (
	"flag"
	"log"
	"net/http"
	"os/exec"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	token := flag.String("token", "xxxx-xxxx-xxxx-xxxx-xxxx-xxxx", "set token")
	address := flag.String("l", "0.0.0.0:8080", "set server listen address and port")
	help := flag.Bool("h", false, "this help")
	script := "./script.sh"
	flag.Parse()

	if *help {
		flag.Usage()
		return
	}

	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		uToken := r.URL.Query().Get("token")
		if *token != uToken {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		cmd := exec.Command(script)
		stdoutStderr, err := cmd.CombinedOutput()
		if err != nil {
			log.Panic(err)
		}
		w.Write(stdoutStderr)
	})

	log.Println("Listen address:", *address)
	log.Fatal(http.ListenAndServe(*address, r))
}
