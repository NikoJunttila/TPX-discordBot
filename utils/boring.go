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
func CheckPromotionDemotion(previousTier, previousRank, currentTier, currentRank string) (string, bool) {
	tierOrder := []string{"IRON", "BRONZE", "SILVER", "GOLD", "PLATINUM", "DIAMOND", "MASTER", "GRANDMASTER", "CHALLENGER"}
	rankOrder := []string{"IV", "III", "II", "I"}

	previousTierIndex := indexOf(tierOrder, previousTier)
	currentTierIndex := indexOf(tierOrder, currentTier)
	demote := "https://media1.tenor.com/m/1CzTPzU3jCwAAAAd/f1-fernando-alonso.gif"
	promote := "https://media1.tenor.com/m/1dopRQonB-IAAAAd/fernando-alonso-fernando.gif"
	if currentTierIndex > previousTierIndex {
		return fmt.Sprintf("Promoted to a higher tier to %s %s \n %s", currentTier, currentRank, promote), true
	} else if currentTier == previousTier {
		previousRankIndex := indexOf(rankOrder, previousRank)
		currentRankIndex := indexOf(rankOrder, currentRank)

		if currentRankIndex > previousRankIndex {
			return fmt.Sprintf("Promoted within the same tier to %s %s \n %s", currentTier, currentRank, promote), true
		} else if currentRankIndex < previousRankIndex {
			return fmt.Sprintf("Demoted within the same tier to %s %s \n %s", currentTier, currentRank, demote), true
		}
	} else {
		return fmt.Sprintf("Demoted to a lower tier to %s %s \n %s", currentTier, currentRank, demote), true
	}
	return "No change", false
}

type tierLP struct {
	tier string
	lp   int
}
type rankLP struct {
	rank string
	lp   int
}

func RankToLP(tier, rank string, lp int) int {
	tierOrder := []tierLP{{tier: "IRON", lp: 0}, {tier: "BRONZE", lp: 400}, {tier: "SILVER", lp: 800}, {tier: "GOLD", lp: 1200}, {tier: "PLATINUM", lp: 1600}, {tier: "EMERALD", lp: 2000}, {tier: "DIAMOND", lp: 2000}, {tier: "MASTER", lp: 2400}, {tier: "GRANDMASTER", lp: 2400}, {tier: "CHALLENGER", lp: 2400}}
	rankOrder := []rankLP{{rank: "IV", lp: 0}, {rank: "III", lp: 100}, {rank: "II", lp: 200}, {rank: "I", lp: 300}}
	lpTotal := 0
	for _, t := range tierOrder {
		if tier == t.tier {
			lpTotal += t.lp
		}
	}
	for _, t := range rankOrder {
		if rank == t.rank {
			lpTotal += t.lp
		}
	}
	lpTotal += lp
	return lpTotal

}
func indexOf(slice []string, item string) int {
	for i, v := range slice {
		if v == item {
			return i
		}
	}
	return -1
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
