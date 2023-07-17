package main

import (
	// Local:
	"web/db"
	"web/hash"
	"web/scrape"
	"web/thread"
	"web/tokens"

	//"web/gpt"

	// Std:
	"fmt"
	"log"
	"time"

	// External:
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
)

func ReturnEmail(token string) string {
  decoded := tokens.Parse(token, "token")
  email := string(decoded)
  return email
}

func ReturnQuizes(email string) ([]interface{} ,[][]map[string]interface{}) {
  data := db.Read(bson.M{
      "type": "quiz",
  }, email)
  var quizzes [][]map[string]interface{}
  var titles []interface{}
  for i := range data {
    titles = append(titles, data[i]["quiz"])
    var temp []map[string]interface{}
    array := data[i]["title"].(primitive.A)
        for i := range array {
      temp = append(temp, array[i].(primitive.M))
    }
    quizzes = append(quizzes, temp)
  } 
  return titles, quizzes
}

func main() {
  engine := html.New("./frontend/dist", ".html")

	app := fiber.New(fiber.Config{
		Views: engine, //set as render engine
	})
  
  app.Use(func(c *fiber.Ctx) error {
    return c.Next()
  })

  app.Static("/", "./frontend/dist")

  app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", nil) })
 

  app.Get("/register", func(c *fiber.Ctx) error {
    return c.Render("register", nil)
  })

  app.Get("/login", func(c *fiber.Ctx) error {
    if c.Cookies("user") != ""  {
      return c.Redirect("/user")
    } else { return c.Render("login", nil)
    }
  })

  
  app.Get("/register", func(c *fiber.Ctx) error {
    return c.Render("register", nil)
  })


  app.Post("/register", func(c *fiber.Ctx) error {
    var email string = c.FormValue("email")
    var pass string = c.FormValue("pass")
    hashpass, err := hash.HashPassword(pass)
      if err != nil {
        log.Panic(err)
      }
    token := tokens.Build(email, "token")  

    if db.UserExists(email) {
		  return c.Render("register", nil) 
    } else {
      db.Write(bson.D{
        {Key: "email", Value: email},
        {Key: "password", Value: hashpass},
        {Key: "token", Value: token},
      }, email)
      return c.Redirect("/login")
    }
  })

  user := app.Group("/user", func(c *fiber.Ctx) error {
    if tokens.Parse(c.Cookies("token"), "token") != "" {
      return c.Next()
    } else {
      return c.Redirect("/login")
    }
  })

  user.Get("/", func(c *fiber.Ctx) error {
    var email string = ReturnEmail(c.Cookies("token"))
    return c.Render("user", fiber.Map{
      "Email": email,
    })
  })

  user.Get("/view", func(c *fiber.Ctx) error {
    var email string = ReturnEmail(c.Cookies("token"))
    _, quizzes := ReturnQuizes(email)
    var ids []int 
    for i := 0; i < len(quizzes); i++ {
      ids = append(ids, i)
    } 
    return c.Render("user/view", fiber.Map{
      "Id": ids,
    })
  })

  user.Get("/view/:id", func(c *fiber.Ctx) error {
    var email string = ReturnEmail(c.Cookies("token"))
    titles, quizzes := ReturnQuizes(email)
    id, err := c.ParamsInt("id")
    if err != nil || id > len(quizzes) - 1 || id < 0 {
      return c.SendString("Invalid Id")
    }
    var question []map[string]interface{} = quizzes[id]
    var title interface{} = titles[id]
    var nextId int = id
    var previousId int = id
    var isNext bool
    var isPrevious bool
    if id == len(quizzes) - 1 && len(quizzes) - 1 >= 1  {
      previousId -= 1
      isNext = false
      isPrevious = true
    } else if id == 0 && len(quizzes) > 1 {
      nextId += 1
      isNext = true
      isPrevious = false
    } else if id == 0 && len(quizzes) - 1 == 0 {
      isPrevious = false
      isNext = false 
    } else {
      nextId += 1
      previousId -= 1
      isNext = true
      isPrevious = true
    }
    return c.Render("user/viewid", fiber.Map{
      "Title": title,
      "Question": question,
      "NextId": nextId,
      "PreviousId": previousId,
      "isNext": isNext,
      "isPrevious": isPrevious,
    })
  })

  user.Get("/gen", func(c *fiber.Ctx) error {
    return c.Render("user/generate", nil)
  })

  user.Post("/gen", func(c *fiber.Ctx) error {
    var title string = c.FormValue("title")
    var url string = c.FormValue("url")
    var email string = ReturnEmail(c.Cookies("token"))
    var paragraphs []string = scrape.GetP(url)
    response := thread.MakQuiz(paragraphs)
    db.WriteQuiz(response, title, "quiz", email)
    return c.Render("user/start/generate", fiber.Map{
        "Title": title,
        "Questions": response,
    })
  })
  app.Post("/login", func(c *fiber.Ctx) error {
      var email string = c.FormValue("email")
      var pass string = c.FormValue("pass")

       data := db.Read(bson.M{
        "email": email,
      }, email)
    
      hashed := data[0]["password"]
      str := fmt.Sprintf("%v", hashed)

      if hash.CheckPasswordHash(pass, str) {
        token := data[0]["token"]
        str := fmt.Sprintf("%v", token)
        cookie := new(fiber.Cookie)
        cookie.Name = "token"
        cookie.Value = str
        cookie.Expires = time.Now().Add(24 * time.Hour)
        cookie.HTTPOnly = true
        c.Cookie(cookie)     
       
        return c.Redirect("/user")
      } else {
        return c.SendString("Oopsie daisy")
    }
  })

  app.Listen(":8080")
}
