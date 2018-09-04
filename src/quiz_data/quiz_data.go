package quiz_data

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type Question struct {
	Category         string   `json:"category"`
	Type             string   `json:"type"`
	Difficulty       string   `json:"difficulty"`
	QuestionText     string   `json:"question"`
	CorrectAnswer    string   `json:"correct_answer"`
	IncorrectAnswers []string `json:"incorrect_answers"`
}

type QuizData struct {
	ResponseCode int        `json:"response_code"`
	Questions    []Question `json:"results"`
}

func GetQuiz(cat string, questNum string) (*QuizData) {
	url := fmt.Sprintf("https://opentdb.com/api.php?amount=%s&category=%s&difficulty=easy&type=boolean",
		questNum, cat)
	response, err := http.Get(url)
	fmt.Println(url)
	fmt.Println(response)

	if err != nil {
		fmt.Printf("%s", err)
		return nil
	} else {
		defer response.Body.Close()
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			fmt.Printf("%s", err)
			return nil
		}

		var r QuizData
		err = json.Unmarshal(contents, &r)
		if err != nil {
			fmt.Printf("err was %v", err)
		}

		for i, _ := range r.Questions {
			q := &r.Questions[i]
			q.QuestionText = strings.Replace(q.QuestionText, "&quot;", "\"", -1)
			q.QuestionText = strings.Replace(q.QuestionText, "&#039;", `'`, -1)
		}

		fmt.Println(r)
		return &r
	}
}
