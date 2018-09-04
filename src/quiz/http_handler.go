package quiz

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	. "quiz_data"

	"github.com/nlopes/slack"
)

// InteractionHandler handles interactive message response.
type InteractionHandler struct {
	slackClient       *slack.Client
	VerificationToken string
}

func (h InteractionHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	fmt.Println("handle_interactive_event()")

	if r.Method != http.MethodPost {
		log.Printf("[ERROR] Invalid method: %s", r.Method)
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Printf("[ERROR] Failed to read request body: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	jsonStr, err := url.QueryUnescape(string(buf)[8:])
	println("jsonStr:", jsonStr)
	if err != nil {
		log.Printf("[ERROR] Failed to unespace request body: %s", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	var message slack.AttachmentActionCallback
	if err := json.Unmarshal([]byte(jsonStr), &message); err != nil {
		log.Printf("[ERROR] Failed to decode json message from slack: %s", jsonStr)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Only accept message from slack with valid token
	if message.Token != h.VerificationToken {
		log.Printf("[ERROR] Invalid token: %s", message.Token)
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Only accept user reaction on the current and relevant action
	if UsersData.UsersCount() == 0 ||  message.CallbackID != UsersData.GetCallbackID(message.User.ID) {
		fmt.Println("not_allowed_user_action")
		w.WriteHeader(http.StatusConflict)
		return
	}

	uid := message.User.ID
	action := message.Actions[0]

	origMsg := message.OriginalMessage
	var attachment slack.Attachment
	var cid string

	switch action.Name {
	case ActionSelectCategory:
		cat := action.SelectedOptions[0].Value
		UsersData.SetCategory(uid, cat)
		return

	case ActionSelectQuestionNum:
		qn := action.SelectedOptions[0].Value
		UsersData.SetQuestionNum(uid, qn)
		return

	case ActionStartQuiz:
		if UsersData.GetCategory(uid) != "" && UsersData.GetQuestionNum(uid) != ""{
			qd := GetQuiz(UsersData.GetCategory(uid), UsersData.GetQuestionNum(uid))
			UsersData.SetQuizData(uid, *qd)

			q := UsersData.NextQuestion(uid)
			attachment, cid = ComposeQuizQuestion(q)
			UsersData.SetCallbackID(uid, cid)
		} else {
			return
		}

	case ActionAnswerTrue, ActionAnswerFalse:
		// Process user answer to the current question
		q := UsersData.CurrentQuestion(uid)
		if (q.CorrectAnswer == "True" && action.Name == ActionAnswerTrue) ||
			(q.CorrectAnswer == "False" && action.Name == ActionAnswerFalse) {
			UsersData.IncCorrectAnswers(uid)
		}

		q = UsersData.NextQuestion(uid)
		if q != nil {
			// Ask next question
			attachment, cid = ComposeQuizQuestion(q)
			UsersData.SetCallbackID(uid, cid)
		} else {
			// Finish the quiz
			attachment, _ = ComposeQuizEnd(
				UsersData.GetCorrectAnswersCount(uid),
				UsersData.GetTotalQuestionsCount(uid))
			UsersData.Remove(uid)
		}
	}

	origMsg.Attachments = []slack.Attachment{
		attachment,
	}
	w.Header().Add("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(&origMsg)
}
