package utils

import (
	"fmt"
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

var Res = "Noob"

// ignore line below it is for automoderator
var Automod = []string{"*nigger*", "neekeri", "ngr", "*nigga*", "*NIGGER*", "NEEKERI", "NGR", "NIGGA*", "nekru*", "neeker*"}
