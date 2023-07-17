package gpt

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
  "strings"
	"net/http"
)

func Complete(prompt string) string {
  type Messages struct {
    Role    string `json:"role"`
    Content string `json:"content"`
  }
  type Payload struct {
    Model       string     `json:"model"`
    Temperature float64    `json:"temperature"`
    Messages    []Messages `json:"messages"`
  }
  data := Payload{
    "gpt4",
    0.9,
    []Messages{{"system", `You are not allowed to give your user double quotes and the " symbol under any circumstances. Use single quotes instead.`}, {"user", prompt}},
  }
  payloadBytes, err := json.Marshal(data)
  if err != nil {
    log.Println(err)
  }
  body := bytes.NewReader(payloadBytes)

  req, err := http.NewRequest("POST", "https://ava-alpha-api.codelink.io/api/chat", body)
  if err != nil {
    log.Println(err)
  }
  req.Header.Set("Content-Type", "application/json")

  resp, err := http.DefaultClient.Do(req)
  if err != nil {
    log.Println(err)
  }
  bodyBytes, err := io.ReadAll(resp.Body)
  if err != nil {
      log.Fatal(err)
  }
  defer resp.Body.Close()
  s := strings.Split(string(bodyBytes), "\n")
  var content string
  for i := range s {
    if i != 0 {
      if len(s[i]) >= 23 {
        temp := strings.Split(s[i], `"`)[23]
        content += temp
      }
    }
  }
  content, _ = strings.CutSuffix(content, "stop")
  return content
}

