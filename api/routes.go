package api

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"time"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/hiteshjoshi/manners"
	"github.com/hiteshjoshi/negroni"
	"github.com/hiteshjoshi/pandat/clock"
)

type Engine struct {
	Router *mux.Router

	Server *manners.GracefulServer

	Clock *clock.Clock
}

func NewRouter() *Engine {

	e := Engine{
		Router: mux.NewRouter().StrictSlash(true),
	}

	return &e
}

func (e *Engine) Stop() {
	e.Server.Close()
}
func (e *Engine) Run(port string) {
	//SetFormat
	n := negroni.New()
	logger := negroni.NewLogger()

	logger.SetAppName("[pandat]")
	logger.SetFormat("{{.AppName}} > {{.StartTime}} | {{.Status}} | \t {{.Duration}} | {{.Hostname}} | {{.Method}} {{.Path}} \n")

	n.Use(logger)

	n.UseHandler(e.Router)

	h := handlers.CORS()(handlers.CompressHandler(handlers.RecoveryHandler()(n)))

	e.Server = manners.NewWithServer(&http.Server{
		Handler: h,
		Addr:    port,
		// Good practice: enforce timeouts for servers you create!
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	})

	//Handling server Interrupt events
	s := make(chan os.Signal, 2)
	signal.Notify(s, os.Interrupt)
	go func() {
		for range s {
			//stop cron server
			fmt.Println("\n\n\nPlease wait....")
			fmt.Println("\nStopping API server!")
			e.Server.Close()
			os.Exit(0)
		}
	}()

	log.Fatal(e.Server.ListenAndServe())
}

func (e *Engine) Get(path string, c http.HandlerFunc) {
	e.Router.
		Methods("GET").
		Path(path).
		Handler(c)
}

func (e *Engine) Post(path string, c http.HandlerFunc) {
	e.Router.
		Methods("POST").
		Path(path).
		Handler(c)
}

func (e *Engine) Put(path string, c http.HandlerFunc) {
	e.Router.
		Methods("PUT").
		Path(path).
		Handler(c)
}

func (e *Engine) Delete(path string, c http.HandlerFunc) {
	e.Router.
		Methods("DELETE").
		Path(path).
		Handler(c)
}

func reportToSentry(error interface{}) {
	// write code here to report error to Sentry
	panic(error)
}
