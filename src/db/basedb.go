package db

import . "quiz_data"

type BaseDB interface {
	UsersCount() int
	Init(uid string)
	IncCorrectAnswers(uid string)
	CurrentQuestion(uid string) *Question
	NextQuestion(uid string) *Question
	Remove(uid string)
	GetCorrectAnswersCount(uid string) int
	GetTotalQuestionsCount(uid string) int
	SetCallbackID(uid string, cid string)
	GetCallbackID(uid string) string
	SetCategory(uid string, cat string)
	GetCategory(uid string) string
	SetQuizData(uid string, qd QuizData)
	SetQuestionNum(uid string, qn string)
	GetQuestionNum(uid string) string
}
