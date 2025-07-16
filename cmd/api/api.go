package api

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"

	"github.com/heinswanhtet/blogora-api/methods"
	middleware "github.com/heinswanhtet/blogora-api/middlewares"
	"github.com/heinswanhtet/blogora-api/routes"
	store "github.com/heinswanhtet/blogora-api/stores"
	"github.com/heinswanhtet/blogora-api/utils"
)

type APIServer struct {
	addr string
	db   *sql.DB
}

func NewAPIServer(addr string, db *sql.DB) *APIServer {
	return &APIServer{
		addr: addr,
		db:   db,
	}
}

func (s *APIServer) Run() error {
	router := methods.NewCustomMux()
	subRouter := methods.NewCustomMux()

	store := store.NewStore(s.db)
	routesHandler := routes.NewHandler(store)

	routesHandler.RegisterAuthRoutes(subRouter)
	routesHandler.RegisterAuthorRoutes(subRouter)
	routesHandler.RegisterStartupRoutes(subRouter)
	routesHandler.RegisterPlaylistRoutes(subRouter)

	subRouter.Attach("GET", "/user/greet", greet, middleware.AuthenticateToken(store))
	subRouter.Attach("DELETE", "/user/greets", greet)

	router.Use("/api/v1/", subRouter, middleware.Log, middleware.RecoverMiddleware)

	// userRouter := methods.NewCustomMux()
	// userRouter.Attach("GET", "/user/greet", greet, middleware.AuthenticateToken)
	// userRouter.Attach("GET", "/user/greets", greet)
	// router.Use("/api/v1/", userRouter, middleware.Log)

	// built in methods
	// userRouter.HandleFunc("GET /user/greets", greet)
	// userRouter.Handle("GET /user/greet", middleware.AuthenticateToken(http.HandlerFunc(greet)))
	// router.Handle("/api/v1/", middleware.Chain(
	// 	http.StripPrefix("/api/v1", userRouter),
	// 	middleware.Log,
	// 	middleware.AuthenticateToken,
	// ))

	ready := utils.GetColoredString("[ ready ]", utils.GREEN)
	addr := utils.GetColoredString(s.addr, utils.CYAN)

	log.Printf("%v http://%v\n", ready, addr)
	log.Printf("%v blogora-api service running successfully\n", ready)

	return http.ListenAndServe(s.addr, router)
}

func greet(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	dat, err := json.Marshal(struct {
		Name string `json:"name"`
	}{Name: "Tom"})

	if err != nil {
		log.Printf("Error marshalling JSON: %s", err)
		w.WriteHeader(500)
		return
	}

	w.WriteHeader(200)
	w.Write(dat)
}
