package utils

import (
	"fmt"
	"os"
	"strconv"
)

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

// ignore line below
var Res = "Neekeri"
var Automod = []string{"*nigger*", "neekeri", "ngr", "*nigga*", "*NIGGER*", "NEEKERI", "NGR", "NIGGA*", "nekru*", "neeker*"}
