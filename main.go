package main

import (
	// Local:
	"web/db"
	"web/hash"
	"web/scrape"

	// Std:
	"encoding/base64"
	"fmt"
	"log"
	"time"

	// External:
	"go.mongodb.org/mongo-driver/bson"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
) 

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

  app.Get("/user", func(c *fiber.Ctx) error {
    if c.Cookies("user") != "" {
      token := c.Cookies("user")
      decoded, err := base64.StdEncoding.DecodeString(token)
	    if err != nil {
        log.Fatal(err)
	    }
      email := string(decoded)
      log.Println(email)
      return c.Render("user", nil)
    } else {
      return c.Redirect("/login")
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
    token := base64.StdEncoding.EncodeToString([]byte(email))   

    if db.UserExists(email) {
		  return c.Render("register", nil) 
    } else {
      db.Write(bson.D{
        {Key: "email", Value: email},
        {Key: "password", Value: hashpass},
        {Key: "token", Value: token},
      }, "quiz", email)
      return c.Redirect("/login")
    }
  })

  app.Get("/gen", func(c *fiber.Ctx) error {
    return c.Render("generate", nil)
  })

  app.Post("/gen", func(c *fiber.Ctx) error {
    var title string = c.FormValue("title")
    var url string = c.FormValue("url")
    log.Println(title)
    var paragraphs []string = scrape.GetP(url)
    log.Println(paragraphs[2])
    return c.Render("generate", nil)
  })

  app.Post("/login", func(c *fiber.Ctx) error {
      var email string = c.FormValue("email")
      var pass string = c.FormValue("pass")

       data := db.Read(bson.M{
        "email": email,
      },  "quiz", email)
    
      hashed := data[0]["password"]
      str := fmt.Sprintf("%v", hashed)

      if hash.CheckPasswordHash(pass, str) {
        token := data[0]["token"]
        str := fmt.Sprintf("%v", token)
        cookie := new(fiber.Cookie)
        cookie.Name = "user"
        cookie.Value = str
        cookie.Expires = time.Now().Add(30 * time.Minute)
        c.Cookie(cookie)     
       
        return c.Redirect("/user")
      } else {
        return c.SendString("Oopsie daisy")
    }
  })

  app.Listen(":8080")
}
