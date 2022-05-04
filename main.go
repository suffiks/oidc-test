package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/coreos/go-oidc/v3/oidc"
	"github.com/davecgh/go-spew/spew"
)

func main() {
	ctx := context.Background()
	provider, err := oidc.NewProvider(ctx, "https://token.actions.githubusercontent.com")
	if err != nil {
		panic(err)
	}

	http.HandleFunc("/token", login(provider))

	fmt.Println("Listening on :8084")
	log.Fatal(http.ListenAndServe(":8084", nil))
}

type ghClaims struct {
	Ref               string `json:"ref"`
	Repository        string `json:"repository"`
	RepositoryID      string `json:"repository_id"`
	RepositoryOwner   string `json:"repository_owner"`
	RepositoryOwnerID string `json:"repository_owner_id"`
	RunID             string `json:"run_id"`
	RunNumber         string `json:"run_number"`
	RunAttempt        string `json:"run_attempt"`
	Actor             string `json:"actor"`
	ActorID           string `json:"actor_id"`
	Workflow          string `json:"workflow"`
	HeadRef           string `json:"head_ref"`
	BaseRef           string `json:"base_ref"`
	EventName         string `json:"event_name"`
	RefType           string `json:"ref_type"`
	Environment       string `json:"environment"`
	JobWorkflowRef    string `json:"job_workflow_ref"`
}

func login(provider *oidc.Provider) http.HandlerFunc {
	verifier := provider.Verifier(&oidc.Config{
		SkipClientIDCheck: true,
	})

	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("X-Token")
		if token == "" {
			http.Error(w, "missing token", http.StatusBadRequest)
			return
		}

		idToken, err := verifier.Verify(r.Context(), token)
		if err != nil {
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		claims := &ghClaims{}
		if err := idToken.Claims(claims); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		fmt.Println("Audience:", idToken.Audience)
		spew.Dump(claims)
	}
}
