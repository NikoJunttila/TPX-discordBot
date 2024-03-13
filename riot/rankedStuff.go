package riot

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func getSummonerInfoByPuuID(puuID string, apiKey string, region string) (SummonerInfo, error) {
	var summonerInfo SummonerInfo
	// Construct the API endpoint URL
	apiEndpoint := fmt.Sprintf("https://%s.api.riotgames.com/lol/summoner/v4/summoners/by-puuid/%s", region, puuID)

	// Create a new HTTP request with the specified method, URL, and optional body
	req, err := http.NewRequest("GET", apiEndpoint, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return summonerInfo, err
	}
	req.Header.Set("X-Riot-Token", apiKey)
	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making HTTP request:", err)
		return summonerInfo, err
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		fmt.Printf("Error: %s\n", response.Status)
		return summonerInfo, err
	}
	body, err := io.ReadAll(response.Body)
	if err != nil {
		fmt.Println("Error copying response body:", err)
		return summonerInfo, err
	}
	// Create an instance of the SummonerInfo struct

	// Unmarshal the JSON data into the struct
	err = json.Unmarshal(body, &summonerInfo)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return summonerInfo, err
	}

	return summonerInfo, nil
}
func getRankedStats(id string, apiKey string, region string) (LeagueEntry, error) {
	var leagueEntries []LeagueEntry
	apiEndpoint2 := fmt.Sprintf("https://%s.api.riotgames.com/lol/league/v4/entries/by-summoner/%s", region, id)
	client := &http.Client{}
	req2, err := http.NewRequest("GET", apiEndpoint2, nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return leagueEntries[0], err
	}
	req2.Header.Set("X-Riot-Token", apiKey)
	response2, err := client.Do(req2)
	if err != nil {
		fmt.Println("Error making HTTP request:", err)
		return leagueEntries[0], err
	}
	defer response2.Body.Close()
	if response2.StatusCode != http.StatusOK {
		fmt.Printf("Error2: %s\n", response2.Status)
		return leagueEntries[0], err
	}
	body2, err := io.ReadAll(response2.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return leagueEntries[0], err
	}
	// Create a slice of LeagueEntry to unmarshal the JSON array

	// Unmarshal the JSON data into the slice
	err = json.Unmarshal(body2, &leagueEntries)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return leagueEntries[0], err
	}
	if len(leagueEntries) <= 1 {
		return leagueEntries[0], err
	}
	if leagueEntries[0].QueueType == "RANKED_SOLO_5x5" {
		return leagueEntries[0], nil
	} else {
		return leagueEntries[1], nil
	}
}
func RankedStats(puuID string, apiKey string, region string) (LeagueEntry, error) {
	var stats LeagueEntry
	summonerInf, err := getSummonerInfoByPuuID(puuID, apiKey, region)
	if err != nil {
		fmt.Println(err)
		return stats, err
	}
	id := summonerInf.ID
	stats, err = getRankedStats(id, apiKey, region)
	if err != nil {
		fmt.Println(err)
		return stats, err
	}
	fmt.Println(stats)
	return stats, nil
}
