package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	"github.com/rauldotgit/pricemage/handlers"
	"github.com/rauldotgit/pricemage/utils"
)

func main() {
	fmt.Println("Initializing ...")

	// loading env variable
	godotenv.Load()
	portString := os.Getenv("PORT")
	if portString == "" {
		log.Fatal("Unable to load PORT env variable.")
	}

	// create chi router
	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*", "http://localhost:8080"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	// new version router
	v1Router := chi.NewRouter()

	// add request handlers
	v1Router.Get("/ok", handlers.HandleOK)

	// add the sub router into the top router
	router.Mount("/v1", v1Router)

	// local execution for testing
	go utils.LocalExec()

	server := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}

	fmt.Println(`                                                                                                                                             

 ______   ______     __     ______     ______     __    __     ______     ______     ______    
/\  == \ /\  == \   /\ \   /\  ___\   /\  ___\   /\ "-./  \   /\  __ \   /\  ___\   /\  ___\   
\ \  _-/ \ \  __<   \ \ \  \ \ \____  \ \  __\   \ \ \-./\ \  \ \  __ \  \ \ \__ \  \ \  __\   
 \ \_\    \ \_\ \_\  \ \_\  \ \_____\  \ \_____\  \ \_\ \ \_\  \ \_\ \_\  \ \_____\  \ \_____\ 
  \/_/     \/_/ /_/   \/_/   \/_____/   \/_____/   \/_/  \/_/   \/_/\/_/   \/_____/   \/_____/ 
                                                                                               

	`)
	fmt.Println("Starting on Port:", portString)

	err := server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
