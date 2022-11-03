package main

import (
	"awesomego/src/common"
	"encoding/json"
	"fmt"
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

type ResponseBody struct {
	Response []Name `json:"response"`
}
type ResponseB struct {
	Response map[string][]Name
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
	router.GET("/test", test)
	router.GET("/testarray", testarray)
	router.GET("/name", getName)
	router.GET("/joke", getJoke)
	router.GET("/status", getStatus)
	router.POST("/post", testPost)

	router.GET("/welcome", func(c *gin.Context) {

		firstname := c.DefaultQuery("firstname", "Guest")
		lastname := c.Query("lastname") // shortcut for c.Request.URL.Query().Get("lastname")
		if lastname == "" || len(lastname) == 0 {
			fmt.Print("last name is empty")
		}
		c.String(http.StatusOK, "Hello %s %s", firstname, lastname)
	})

	router.Run("0.0.0.0:8083")
}

func testarray(c *gin.Context) {
	var names []Name
	names = append(names, Name{FirstName: "syn1", LastName: "sun1"})
	names = append(names, Name{FirstName: "syn2", LastName: "sun2"})
	initMap := make(map[string][]Name)
	initMap["roles"] = names
	fmt.Println("names is ", names)
	response := ResponseB{Response: initMap}
	c.JSON(http.StatusOK, response)
}

func test(c *gin.Context) {
	data := make(map[string]string)
	data["a"] = "a"
	data["b"] = "b"
	//s := "this is a test"

	filter := c.Request.Form
	f := make(map[string]interface{})
	log.Println(filter)
	// convert map[string][]string to map[string]interface{}
	ConvertFilterMaps(filter, f)
	log.Println(filter)
	log.Println(f)
	var names []Name
	names = append(names, Name{FirstName: "syn1", LastName: "sun1"})
	names = append(names, Name{FirstName: "syn2", LastName: "sun2"})
	response := ResponseBody{Response: names}
	fmt.Println("names is ", names)
	c.JSON(http.StatusOK, response)
}
func testPost(c *gin.Context) {
	var testName Name
	err := c.ShouldBindJSON(&testName)
	if err != nil {
		log.Println(err)
	} else {
		log.Println("first name is ", testName.FirstName, " and last name is ", testName.LastName)
	}
	c.JSON(http.StatusOK, testName)
}

// convert map[string][]string to map[string]interface{}
func ConvertFilterMaps(inFilter map[string][]string, outFilter map[string]interface{}) {
	for i, j := range inFilter {
		if j[0] == "true" {
			outFilter[i] = true
		} else if j[0] == "false" {
			outFilter[i] = false
		} else {
			outFilter[i] = j[0] // using only the first value
		}
	}
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
