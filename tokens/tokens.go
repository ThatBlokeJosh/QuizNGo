package tokens

import (
	"time"
  "fmt"
	"github.com/golang-jwt/jwt/v5"
	"github.com/joho/godotenv"
)
var hmacSecret []byte

func init() {
  env, _ := godotenv.Read(".env")
  hmacSecret = []byte(env["SECRET"])
}

func Build(data string, dataName string) string {
  token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
    dataName: data,
    "nbf": time.Date(2015, 10, 10, 12, 0, 0, 0, time.UTC).Unix(),
  })
  tokenString, _ := token.SignedString(hmacSecret)
  return tokenString
}


func Parse (tokenString string, tokenName string) string {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}		
		return hmacSecret, nil
	})
  if err != nil {
    return ""
  }
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return fmt.Sprintf("%v", claims[tokenName])
	} else {
    return ""
  }
}

