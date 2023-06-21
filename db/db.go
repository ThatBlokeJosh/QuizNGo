package db

import (
	"context"
	"fmt"
	"log"
  "github.com/joho/godotenv"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Uri() string {
  envFile, _ := godotenv.Read(".env")
  return envFile["MONGO_URI"] 
}

type Question struct {
  title string
  options [4]string
  answer string
  number int
}

func makeQuizBson(q []Question, title string) bson.M {
  var quiz bson.M = bson.M{
    title: nil, 
  }
  var array []bson.M
  for i := 0; i <= len(q) - 1; i++ {
    array = append(array, bson.M{
      "title": q[i].title, "optionA": q[i].options[i], "optionB": q[i].options[i], "optionC": q[i].options[i], "optionD": q[i].options[i], "answer": q[i].answer})
  }
  quiz["quiz"] = array
  return quiz
}

func Connect(dbName string) (*mongo.Database, context.Context) {
  // Connecting
  client, err := mongo.NewClient(options.Client().ApplyURI(Uri()))
  if err != nil {
    log.Fatal(err)
  }

  ctx := context.Background()
  err = client.Connect(ctx) 
  if err != nil {
    log.Fatal(err)
  }
  
  db := client.Database(dbName)

  return db, ctx
}
  
func WriteQuiz(qSlice []Question, title string, dbName string, colName string) {
  db, ctx := Connect(dbName)
  collection := db.Collection(colName)

  //var q1 Question = Question{"test", [4]string{"test", "test", "test", "test"}, "1", 1}
  //var q2 Question = Question{"works", [4]string{"test", "test", "test", "test"}, "1", 1}

  var quiz bson.M = makeQuizBson(qSlice, title)
  // Inserting (to insert multiple use InsertMany)
  result, err := collection.InsertMany(ctx, []interface{}{
    quiz,
  })
  if err != nil {
    log.Fatal(err)
  }
  fmt.Println("Result: ", result)
}

func Write(data bson.D, dbName string, colName string) {
  db, ctx := Connect(dbName)
  collection := db.Collection(colName)

  result, err := collection.InsertMany(ctx, []interface{}{
    data,
  })
  if err != nil {
    log.Fatal(err)
  }
  fmt.Println("Result: ", result)
}

func Read(filter bson.M, dbName string, colName string) []bson.M {
  db, ctx := Connect(dbName)
  collection := db.Collection(colName)

  filteredCursor, err := collection.Find(ctx, filter)
  if err != nil {
    log.Fatal(err)
  }
  var filtered []bson.M
  if err = filteredCursor.All(ctx, &filtered); err != nil {
    log.Fatal(err)
  }
  return filtered
}

func UserExists(email string) bool {
  data := Read(bson.M{
    "email": email,
  }, "quiz", email)

  if data != nil {
    return true
  } else {
    return false
  }
}



