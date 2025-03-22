package router

import (
	"net/http"
	"sibir2025/service/app_logic/handlers"
)

func Setup_routes() {
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	http.HandleFunc("/", handlers.Index)
	http.HandleFunc("GET /sign_in", handlers.Sign_in)
	http.HandleFunc("GET /sign_up", handlers.Sign_up)
	http.HandleFunc("POST /authorize", handlers.Authorize)
	http.HandleFunc("POST /authorize_gosuslugi", handlers.AuthGosuslugi)
	http.HandleFunc("POST /register", handlers.Register)
	http.HandleFunc("POST /logout", handlers.Logout)
	http.HandleFunc("GET /sus_uploader", handlers.Sus_uploader)
	http.HandleFunc("POST /create_sus", handlers.Create_sus)
	http.HandleFunc("GET /sus_browser", handlers.Sus_browser)
	http.HandleFunc("POST /create_claim", handlers.Create_claim)
	http.HandleFunc("POST /accept_claim", handlers.Accept_claim)
	http.HandleFunc("GET /claim_browser", handlers.Claim_browser)
	http.HandleFunc("GET /claim_manager", handlers.Claim_manager)
	http.HandleFunc("GET /sus/", handlers.Sus_profile)
	http.HandleFunc("POST /download_profile", handlers.Download_profile)
}
