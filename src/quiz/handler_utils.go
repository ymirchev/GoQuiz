package quiz

import (
	"db"
	"fmt"
	"github.com/nlopes/slack"
	"math/rand"
	"quiz_data"
)

const (
	ActionAnswerTrue        = "AnswerTrue"
	ActionAnswerFalse       = "AnswerFalse"
	ActionStartQuiz         = "StartQuiz"
	ActionSelectCategory    = "SelectCategory"
	ActionSelectQuestionNum = "SelectQuestionNum"
)

var UsersData db.BaseDB = &db.Users{}

func ComposeQuizQuestion(q *quiz_data.Question) (slack.Attachment, string) {
	cid := generateCallbackID()
	return slack.Attachment{
		Text:       q.QuestionText,
		Color:      "#41f4a0",
		CallbackID: cid,
		Actions: []slack.AttachmentAction{
			{
				Name:  ActionAnswerTrue,
				Text:  "True",
				Type:  "button",
				Style: "primary",
			},
			{
				Name:  ActionAnswerFalse,
				Text:  "False",
				Type:  "button",
				Style: "danger",
			},
		},
	}, cid
}

func ComposeQuizEnd(m, n int) (slack.Attachment, string) {
	cid := generateCallbackID()
	return slack.Attachment{
		Color:      "#f9a41b",
		CallbackID: cid,
		Actions:    []slack.AttachmentAction{}, // empty buttons
		Fields: []slack.AttachmentField{
			{
				Title: "Thank you for taking the quiz!",
				Value: fmt.Sprintf("You answered correctly to %d questions out of %d", m, n),
				Short: false,
			},
		},
	}, cid
}

func ComposeQuizConfig() ([]slack.Attachment, string) {
	cid := generateCallbackID()
	attachment := slack.Attachment{
		Text:       "Make your quiz:",
		Color:      "#f9a41b",
		CallbackID: cid,
		Actions: []slack.AttachmentAction{
			{
				Name: ActionSelectCategory,
				Text: "Category",
				Type: "select",
				Options: []slack.AttachmentActionOption{
					{
						Text:  "General Knowledge",
						Value: "9",
					},
					{
						Text:  "Books",
						Value: "10",
					},
					{
						Text:  "Movies",
						Value: "11",
					},
					{
						Text:  "Music",
						Value: "12",
					},
					{
						Text:  "Computer science",
						Value: "18",
					},
					{
						Text:  "History",
						Value: "23",
					},
				},
			},
			{
				Name: ActionSelectQuestionNum,
				Text: "Number of questions",
				Type: "select",
				Options: []slack.AttachmentActionOption{
					{
						Text:  "1",
						Value: "1",
					},
					{
						Text:  "3",
						Value: "3",
					},
					{
						Text:  "5",
						Value: "5",
					},
					{
						Text:  "10",
						Value: "10",
					},
				},
			},
			{
				Name:  ActionStartQuiz,
				Text:  "Click to go",
				Type:  "button",
				Style: "primary",
			},
		},
	}

	return []slack.Attachment{
		attachment,
	}, cid
}

// Generate a random string of A-Z chars with len = l
func randomString(len int) string {
	bytes := make([]byte, len)
	for i := 0; i < len; i++ {
		bytes[i] = byte(65 + rand.Intn(25)) // A=65 and Z = 65+25
	}
	return string(bytes)
}

func generateCallbackID() string {
	return randomString(20)
}
