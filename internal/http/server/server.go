package server

import (
	"context"
	"crypto/tls"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/torwig/promotions/internal/promotions/app"
)

type Server struct {
	srv *http.Server
	app app.Application
}

func NewServer(application app.Application) Server {
	return Server{app: application}
}

func (s Server) GetPromotion(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	promo, err := s.app.GetPromotion(r.Context(), id)
	if err != nil {
		respondWithError(err, w, r)
		return
	}

	resp := promotionResponse(*promo)
	respond(w, r, resp)
}

func (s *Server) RunHTTPServer(port string) {
	router := chi.NewRouter()
	setMiddlewares(router)

	router.Get("/promotions/{id}", s.GetPromotion)

	log.Println("starting HTTPS server")

    certFile := os.Getenv("SSL_CERTIFICATE")
	keyFile := os.Getenv("SSL_PRIVATE_KEY")

	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		log.Fatalf("unable to load SSL certificate: %+v", err)
	}

	s.srv = &http.Server{
		Addr:         ":" + port,
		Handler:      router,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  120 * time.Second,
		TLSConfig:    &tls.Config{
			Certificates: []tls.Certificate{cert},
		},
	}

	go func() {
		if err := s.srv.ListenAndServeTLS("", ""); err != nil && err != http.ErrServerClosed {
			log.Fatalf("unable to start HTTPS server: %+v", err)
		}
	}()
}

func (s *Server) Shutdown(ctx context.Context) {
	if err := s.srv.Shutdown(ctx); err != nil {
		log.Fatalf("unable to shutdown HTTPS server: %+v", err)
		return
	}

	log.Println("HTTPS server has been stopped")
}

func setMiddlewares(router *chi.Mux) {
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
}
