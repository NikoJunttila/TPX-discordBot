package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"strconv"
)

func WaitForInterruptSignal() {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	fmt.Println("Press Ctrl+C to exit")
	<-stop
}
func GetEnvVariable(key string) string {
	value := os.Getenv(key)
	if value == "" {
		panic(fmt.Sprintln("no env value: ", key))
	}
	return value
}
func CatPic() (string, error) {
	type catType struct {
		Url string `json:"url"`
	}

	var cat []catType
	url := "https://api.thecatapi.com/v1/images/search"

	req, _ := http.NewRequest("GET", url, nil)
	catAPI := GetEnvVariable("catAPI")
	req.Header.Set("x-api-key", catAPI)
	req.Header.Add("accept", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	err = json.Unmarshal(body, &cat)
	if err != nil {
		return "", err
	}
	return cat[0].Url, nil
}
func InsultRes() string {
	type insult struct {
		Insult string `json:"insult"`
	}

	var ins insult
	url := "https://evilinsult.com/generate_insult.php?lang=en&type=json"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("accept", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)
	_ = json.Unmarshal(body, &ins)
	return ins.Insult
}

func IncrementAndWriteToFile(filename string) (int, error) {
	// Read the content of the file
	content, err := os.ReadFile(filename)
	if err != nil {
		return 0, err
	}

	// Parse the content as an integer
	num, err := strconv.Atoi(string(content))
	if err != nil {
		return 0, err
	}

	// Increment the number
	num++

	// Convert the updated number back to string
	updatedContent := strconv.Itoa(num)

	// Write the updated content back to the file
	err = os.WriteFile(filename, []byte(updatedContent), os.ModePerm)
	if err != nil {
		return 0, err
	}

	fmt.Printf("Number updated: %d\n", num)

	return num, nil
}

// ignore line below it is for automoderator
var Automod = []string{"*nigger*", "*neekeri*", "ngr", "*nigga*", "*nekru*"}
