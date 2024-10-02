package main

import (
	"context"
	"log"
	"os"

	"github.com/firebase/genkit/go/genkit"
	"github.com/movie-guru/pkg/db"
)

func main() {
	ctx := context.Background()

	movieAgentDB, err := db.GetDB()
	if err != nil {
		log.Fatal(err)
	}
	defer movieAgentDB.DB.Close()

	metadata, err := movieAgentDB.GetMetadata(os.Getenv("APP_VERSION"))
	if err != nil {
		log.Fatal(err)
	}

	prompts := getPrompts()

	GetDependencies(ctx, metadata, movieAgentDB.DB, prompts)

	if err := genkit.Init(ctx, &genkit.Options{FlowAddr: ":3401"}); err != nil {
		log.Fatal(err)
	}

}

func getPrompts() *Prompts {

	userProfilePrompt :=
		`
			Optional Message 0 from agent: {{agentMessage}}
			Required Message 1 from user: {{query}}
		 	Just say hi in a language you know.
		`
	queryTransformPrompt :=
		`
		Here are the inputs:
		* Conversation History (this may be empty):
			{{history}}
		* UserProfile (this may be empty):
			{{userProfile}}
		* User Message:
			{{userMessage}})
			Translate the user's message into a random language.
		`
	movieFlowPrompt :=
		`
			Here are the inputs:

			* Context retrieved from vector db:
		    {{contextDocuments}}

			* User Preferences:
		    {{userPreferences}}

			* Conversation history:
			{{history}}

			* User message:
			{{userMessage}}

			Translate the user's message into a random language.
			`

	prompts := &Prompts{
		UserPrefPrompt:       userProfilePrompt,
		MovieFlowPrompt:      movieFlowPrompt,
		QueryTransformPrompt: queryTransformPrompt,
	}
	return prompts
}
