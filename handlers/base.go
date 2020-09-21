package handlers

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/dghubble/gologin"
	"github.com/dghubble/gologin/github"
	"github.com/joho/godotenv"
	"golang.org/x/oauth2"
	githubOAuth2 "golang.org/x/oauth2/github"
)

//New func
func New() http.Handler {
	godotenv.Load(".env")
	mux := http.NewServeMux()
	// Root
	config := &oauth2.Config{
		ClientID:     os.Getenv("CLIENTID"),
		ClientSecret: os.Getenv("CLIENTSECRET"),
		RedirectURL:  "http://localhost:8000/oauth/callback",
		Endpoint:     githubOAuth2.Endpoint,
	}
	mux.Handle("/", http.FileServer(http.Dir("templates/")))

	stateConfig := gologin.DebugOnlyCookieConfig
	mux.Handle("/login", github.StateHandler(stateConfig, github.LoginHandler(config, nil)))
	mux.Handle("/oauth/callback", github.StateHandler(stateConfig, github.CallbackHandler(config, issueSession(), nil)))
	mux.Handle("/check", issueSession())
	return mux
}

func issueSession() http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		githubUser, err := github.UserFromContext(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		e := json.NewEncoder(w)
		e.SetIndent("", "  ")
		e.Encode(githubUser)
	}
	return http.HandlerFunc(fn)
}
