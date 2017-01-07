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

func NewEngine() *Engine {

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

			fmt.Println("\nStopping Clock server!")

			if e.Clock.Stop() != nil {
				fmt.Println("\nFailed to remove clock server from redis!")
				fmt.Println("\nClock server id : ", e.Clock.ID)

			} else {
				fmt.Println("\nStopped Clock server!")
			}

			os.Exit(0)
		}
	}()

	fmt.Println("\nStarting API server on port" + port)
	log.Fatal(e.Server.ListenAndServe())
}

func (e *Engine) Group(path string) *Engine {

	ne := NewEngine()
	ne.Clock = e.Clock
	ne.Router = e.Router.PathPrefix(path).Subrouter()
	ne.Server = e.Server
	ne.Router.StrictSlash(true)
	return ne
}
func (e *Engine) Method(m string, path string, c http.HandlerFunc) {
	e.Router.
		Methods(m).
		Path(path).
		Handler(c)
}
func (e *Engine) Get(path string, c http.HandlerFunc) {

	e.Method("GET", path, c)
}

func (e *Engine) WS(path string, c http.HandlerFunc) {

	e.Router.
		Path(path).
		Handler(c)
}

func (e *Engine) Post(path string, c http.HandlerFunc) {

	e.Method("POST", path, c)
}

func (e *Engine) Put(path string, c http.HandlerFunc) {

	e.Method("PUT", path, c)
}

func (e *Engine) Delete(path string, c http.HandlerFunc) {
	e.Method("DELETE", path, c)
}

func reportToSentry(error interface{}) {
	// write code here to report error to Sentry
	panic(error)
}
