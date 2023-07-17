package thread

import (
	"strings"
	"sync"
	"web/gpt"
)

type Question struct {
  Q string
  A string
  B string
  C string
  D string
  Correct string
}

func ParseResponse(s string) Question {
  response := gpt.Complete("Here is some text: " +  s  + " Make a quiz question from and format the response exactly like so: Question: the question? A) option A B) option B C) option C D) option D Correct: Just the letter of the correct option and nothing else.")
  temp := strings.Split(response, "?")
  if len(temp) != 1 {
    var question string = temp[0]
    _, question, _ = strings.Cut(question, "Question: ")
    temp[1] = strings.ReplaceAll(temp[1], `\n`, "")  
    temp2 := strings.Split(temp[1], "B)")
    var a string = temp2[0]
    _, a, _ = strings.Cut(a, "A)")
    temp3 := strings.Split(temp2[1], "C)")
    var b string = temp3[0]
    temp4 := strings.Split(temp3[1], "D)")
    var c string = temp4[0]
    temp5 := strings.Split(temp4[1], "Correct:")
    var d string = temp5[0]
    var correct string = temp5[1]
    var q Question = Question{question, a, b, c, d, correct}
    return q
  } else {
    return ParseResponse(s)
  }
}

func MakQuiz(slice []string) []Question {
  var wg sync.WaitGroup
  wg.Add(len(slice))
  var quiz []Question
  for _, s := range slice {
      go func(s string) {
          defer wg.Done()
          quiz = append(quiz, ParseResponse(s))
      }(s)
  }

  wg.Wait()

  return quiz
}
