package db

import (
	. "quiz_data"
)

type User struct {
	QuizData QuizData
	CurrentQuestion int
	CorrectAnswers int
	CallbackID string
	Category string
	QuestionsNum string
}

type Users map[string]*User

func (u *Users) UsersCount() int {
	return len(*u)
}

func (u *Users) Init(uid string)  {
	(*u)[uid] = &User{QuizData{},-1, 0, "", "", ""}
}

func (u *Users) IncCorrectAnswers(uid string)  {
	(*u)[uid].CorrectAnswers++
}

func (u *Users) CurrentQuestion(uid string) *Question {
	ui := (*u)[uid]
	return &ui.QuizData.Questions[ui.CurrentQuestion]
}

func (u *Users) NextQuestion(uid string) *Question {
	ui := (*u)[uid]

	if ui.CurrentQuestion < u.GetTotalQuestionsCount(uid) - 1 {
		ui.CurrentQuestion++
		return u.CurrentQuestion(uid)
	} else {
		return nil
	}
}

func (u *Users) Remove(uid string) {
	delete(*u, uid)
}

func (u *Users) GetCorrectAnswersCount(uid string) int {
	return (*u)[uid].CorrectAnswers
}

func (u *Users) GetTotalQuestionsCount(uid string) int {
	return len((*u)[uid].QuizData.Questions)
}

func (u *Users) SetCallbackID(uid string, cid string)  {
	(*u)[uid].CallbackID = cid
}

func (u *Users) GetCallbackID(uid string) string {
	return (*u)[uid].CallbackID
}

func (u *Users) SetCategory(uid string, cat string)  {
	(*u)[uid].Category = cat
}

func (u *Users) GetCategory(uid string) string {
	return (*u)[uid].Category
}

func (u *Users) SetQuizData(uid string, qd QuizData)  {
	(*u)[uid].QuizData = qd
}

func (u *Users) SetQuestionNum(uid string, qn string)  {
	(*u)[uid].QuestionsNum = qn
}

func (u *Users) GetQuestionNum(uid string) string {
	return (*u)[uid].QuestionsNum
}
