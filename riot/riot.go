package riot

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func CheckLastMatch(lastGameID string, puuID string, country string, apiKey string) (string, bool) {
	newGameID, err := GetMatchHistory(puuID, 1, country, apiKey)
	if err != nil {
		fmt.Println(err, "err getting match history")
	}
	if lastGameID != newGameID[0] {
		fmt.Print("\n******************************\n*\n* new match \n*\n*******************************\n")
		return newGameID[0], true
	}
	return lastGameID, false
}
func GetMatchHistory(puuID string, count int, country string, apiKey string) ([]string, error) {
	var matches []string
	matchesEndPoint := fmt.Sprintf("https://%s.api.riotgames.com/lol/match/v5/matches/by-puuid/%s/ids?start=0&count=%d", country, puuID, count)
	matchReq, _ := http.NewRequest("GET", matchesEndPoint, nil)
	matchReq.Header.Set("X-Riot-Token", apiKey)
	client := &http.Client{}
	matchRes, err := client.Do(matchReq)
	if err != nil {
		fmt.Println("Error making HTTP request:", err)
		return matches, err
	}
	defer matchRes.Body.Close()
	if matchRes.StatusCode != http.StatusOK {
		fmt.Printf("Error: %s\n", matchRes.Status)
		return matches, err
	}
	body2, err := io.ReadAll(matchRes.Body)
	if err != nil {
		fmt.Println("Error copying response body:errors", err)
		return matches, err
	}
	// Unmarshal the JSON data into the struct
	err = json.Unmarshal(body2, &matches)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return matches, err
	}
	return matches, nil
}
