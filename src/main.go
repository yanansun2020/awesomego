package main

import (
	"awesomego/src/common"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
)

// album represents data about a record album.
type status struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type Name struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

type JokeDetail struct {
	Id         int      `json:"id"`
	Joke       string   `json:"joke"`
	Categories []string `json:"categories"`
}

type Joke struct {
	JokeType string     `json:"type"`
	Value    JokeDetail `json:"value"`
}

func main() {
	common.HttpClientInit()
	router := gin.Default()
	router.GET("/", getCombinedJoke)
	router.GET("/name", getName)
	router.GET("/joke", getJoke)
	router.GET("/status", getStatus)
	router.Run("0.0.0.0:8080")
}

func getCombinedJoke(c *gin.Context) {
	//step1 : get first name and last name
	name, err := fetchRandomName()
	if err != nil {
		msg := "error while fetching random name"
		log.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, getErrorStatus(msg, http.StatusInternalServerError))
		return
	}
	if name.FirstName == "" || name.LastName == "" {
		msg := "first name or last name is empty"
		log.Println(msg)
		c.IndentedJSON(http.StatusInternalServerError, getErrorStatus(msg, http.StatusInternalServerError))
		return
	}
	//step2 : get a random joke
	firstName := "John"
	lastName := "Doe"
	limitTo := "[nerdy]"
	fullName := "John Doe"
	joke, err := fetchRandomJoke(firstName, lastName, limitTo)
	if err != nil {
		msg := "error while get joke"
		log.Println(err)
		c.IndentedJSON(http.StatusInternalServerError, getErrorStatus(msg, http.StatusInternalServerError))
		return
	}
	if joke.JokeType != "success" {
		//err = errors.New("Get client failed")
		return
	}
	// step 3, build new name
	var sb strings.Builder
	sb.WriteString(name.FirstName)
	sb.WriteString(" ")
	sb.WriteString(name.LastName)
	log.Println("the new full name is ", sb.String())
	jokeValue := joke.Value.Joke
	log.Println("the joke is ", jokeValue)
	jokeValue = strings.Replace(jokeValue, fullName, sb.String(), -1)
	log.Println("after replace, the joke is ", jokeValue)
	c.String(http.StatusOK, jokeValue)
}

func getName(c *gin.Context) {
	name, err := fetchRandomName()
	if err != nil {
		log.Printf("error while get random name, %s", err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}
	if name.FirstName == "" || name.LastName == "" {
		log.Println("first name or last name is empty")
		c.Status(http.StatusInternalServerError)
		return
	}
	c.IndentedJSON(http.StatusOK, name)
}

func getJoke(c *gin.Context) {
	joke, err := fetchRandomJoke("John", "Doe", "[nerdy]")
	if err != nil {
		log.Printf("error while fetchRandomJoke to get random joke, %s", err.Error())
		c.Status(http.StatusInternalServerError)
		return
	}
	if joke.JokeType != "success" {
		log.Println("get joke type is not success")
		return
	}
	c.IndentedJSON(http.StatusOK, joke)
}

func getStatus(c *gin.Context) {
	status := status{
		Code:    200,
		Message: "OK",
	}
	c.IndentedJSON(http.StatusOK, status)
}

func fetchRandomName() (out Name, err error) {
	var req common.HttpReq
	req.Url = common.NameServer
	req.Method = "GET"
	resp, err := common.SendHttpReq(req)
	if err != nil {
		log.Printf("error while SendHttpReq to get name, %s", err.Error())
		return
	}
	if resp.Status != http.StatusOK {
		//err = errors.New("Get failed")
		return
	}
	json.Unmarshal(resp.Body, &out)
	return
}

func fetchRandomJoke(firstName string, lastName string, limitTo string) (out Joke, err error) {
	var req common.HttpReq
	req.Url = common.JokeAPIServer + "?firstName=" + firstName + "&lastName=" + lastName + "&limitTo=" + limitTo
	req.Method = "GET"
	resp, err := common.SendHttpReq(req)
	if err != nil {
		log.Printf("error while SendHttpReq to get joke, %s", err.Error())
		return
	}
	if resp.Status != http.StatusOK {
		//err = errors.New("Get failed")
		return
	}
	json.Unmarshal(resp.Body, &out)
	return
}
func getErrorStatus(message string, code int) status {
	status := status{
		Code:    code,
		Message: message,
	}
	return status
}
