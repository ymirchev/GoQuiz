package main

import (
	"github.com/nlopes/slack"
	"log"
	"net/http"
	. "quiz"
)

func main() {
	var env EnvConfig = Config()

	// Listening slack event and response
	log.Printf("[INFO] Start slack event listening")
	client := slack.New(env.BotToken)
	slackListener := &SlackListener{Client: client}
	go slackListener.ListenAndResponse()

	// Register handler to receive interactive message
	// responses from slack (kicked by user action)
	http.Handle("/interaction", InteractionHandler{
		VerificationToken: env.VerificationToken,
	})

	log.Printf("[INFO] Server listening on :%s", env.Port)
	if err := http.ListenAndServe(":"+env.Port, nil); err != nil {
		log.Printf("[ERROR] %s", err)
		return
	}
}
