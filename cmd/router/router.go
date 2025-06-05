package router

import (
	"fmt"
	"net/http"
	"path/filepath"
	"strings"

	"net/http/pprof"

	project "github.com/TechXTT/tf-api-TORM/cmd/projects"
	votes "github.com/TechXTT/tf-api-TORM/cmd/votes"
	"github.com/TechXTT/tf-api-TORM/torm/models"
	"github.com/rs/cors"

	"github.com/gorilla/mux"
)

type Router struct {
	router *mux.Router
}

func LimitRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		allowedDomains := []string{"tuesfest.bg", "*.tuesfest.bg", "*.vercel.app", "localhost", "localhost"}

		reqDomain := strings.Split(r.Host, ":")[0]

		for _, domain := range allowedDomains {
			if matched, _ := filepath.Match(domain, reqDomain); matched {
				next.ServeHTTP(w, r)
				return
			}
		}

		http.Error(w, "Forbidden", http.StatusForbidden)

		next.ServeHTTP(w, r)
	})
}

func NewRouter() *Router {
	r := mux.NewRouter().PathPrefix("/v1").Subrouter().StrictSlash(true)

	r.Use(mux.CORSMethodMiddleware(r))
	r.Use(LimitRequest)

	return &Router{r}
}

func (r *Router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.router.ServeHTTP(w, req)
}

func (r *Router) GetRouter() *mux.Router {
	return r.router
}

func (r *Router) Init() {
	router := r.GetRouter()
	router.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("Hello World!"))
	}).Methods("GET")
}

func (r *Router) Pprof() {
	router := r.GetRouter()
	debugProf := router.PathPrefix("/debug/pprof").Subrouter()
	debugProf.HandleFunc("/", pprof.Index)
	debugProf.HandleFunc("/cmdline", pprof.Cmdline)
	debugProf.HandleFunc("/symbol", pprof.Symbol)
	debugProf.HandleFunc("/trace", pprof.Trace)
	debugProf.HandleFunc("/profile", pprof.Profile)

	// Manually add support for paths not easily linked as above
	// Hooking this up is actually very convoluted and only a few answers on how to do it
	// https://stackoverflow.com/questions/19591065/profiling-go-web-application-built-with-gorillas-mux-with-net-http-pprof
	debugProf.Handle("/goroutine", pprof.Handler("goroutine"))
	debugProf.Handle("/heap", pprof.Handler("heap"))
	debugProf.Handle("/threadcreate", pprof.Handler("threadcreate"))
	debugProf.Handle("/block", pprof.Handler("block"))
	debugProf.Handle("/vars", http.DefaultServeMux)
}

func (r *Router) Projects() {
	router := r.GetRouter()
	GetReq := router.PathPrefix("/get").Subrouter().StrictSlash(true)
	GetReq.HandleFunc("/projects", func(writer http.ResponseWriter, request *http.Request) {
		client := models.NewClient()
		// call function GetProjects from projects package
		project.GetProjectsTorm(writer, request, client)
	}).Methods("GET")
	GetReq.HandleFunc("/project/{id}", func(writer http.ResponseWriter, request *http.Request) {
		client := models.NewClient()
		// call function GetProject from projects package
		project.GetProjectTorm(writer, request, client)
	}).Methods("GET")
	GetReq.HandleFunc("/projects/{category}", func(writer http.ResponseWriter, request *http.Request) {
		client := models.NewClient()
		// call function GetProjectsByCategory from projects package
		project.GetProjectsByCategoryTorm(writer, request, client)
	}).Methods("GET")
}

func (r *Router) Votes() {
	router := r.GetRouter()
	PostReq := router.PathPrefix("/post").Subrouter().StrictSlash(true)
	UpdateReq := router.PathPrefix("/update").Subrouter().StrictSlash(true)
	PostReq.HandleFunc("/vote", func(writer http.ResponseWriter, request *http.Request) {
		client := models.NewClient()
		// call function PostVote from projects package
		votes.PostVoteTorm(writer, request, client)
	}).Methods("POST")
	UpdateReq.HandleFunc("/verify_vote", func(writer http.ResponseWriter, request *http.Request) {
		client := models.NewClient()
		// call function VerifyVote from projects package
		votes.VerifyVoteTorm(writer, request, client)
	}).Methods("PUT")
}

func (r *Router) Run() {
	// r.Database()
	r.Projects()
	r.Votes()
	r.Init()
	fmt.Println("Routes initialized")

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"https://tuesfest.bg", "https://tuesfest.bg/", "https://*.tuesfest.bg", "https://*.tuesfest.bg/", "https://*.vercel.app", "https://*.vercel.app/", "http://localhost:3000", "http://localhost:8080"},
		AllowCredentials: true,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
	})

	handler := c.Handler(r)

	err := http.ListenAndServe(":8080", handler)
	if err != nil {
		fmt.Println(err)
		return
	}
}
