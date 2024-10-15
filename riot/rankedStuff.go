package riot

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/nikojunttila/discord/utils"
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
		return LeagueEntry{}, err
	}
	req2.Header.Set("X-Riot-Token", apiKey)
	response2, err := client.Do(req2)
	if err != nil {
		fmt.Println("Error making HTTP request:", err)
		return LeagueEntry{}, err
	}
	defer response2.Body.Close()
	if response2.StatusCode != http.StatusOK {
		fmt.Printf("Error2: %s\n", response2.Status)
		return LeagueEntry{}, err
	}
	body2, err := io.ReadAll(response2.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return LeagueEntry{}, err
	}
	// Create a slice of LeagueEntry to unmarshal the JSON array

	// Unmarshal the JSON data into the slice
	err = json.Unmarshal(body2, &leagueEntries)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return LeagueEntry{}, err
	}
	fmt.Println(leagueEntries)
	var noobPlayer LeagueEntry
	if len(leagueEntries) < 1 {
		noobPlayer.Rank = "IV"
		noobPlayer.Tier = "IRON"
		noobPlayer.LeaguePoints = 1
		noobPlayer.Wins = 1
		noobPlayer.Losses = 1
		noobPlayer.SummonerName = "no ranked noob"
		return noobPlayer, nil
	}
	if leagueEntries[0].QueueType == "RANKED_SOLO_5x5" {
		return leagueEntries[0], nil
	} else if len(leagueEntries) > 1 {
		return leagueEntries[1], nil
	} else {
		noobPlayer.Rank = "IV"
		noobPlayer.Tier = "IRON"
		noobPlayer.LeaguePoints = 1
		noobPlayer.Wins = 1
		noobPlayer.Losses = 1
		noobPlayer.SummonerName = "only flex noob"
		return noobPlayer, nil
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
	return stats, nil
}

func LiveGamePlayersPuuIDS(apiKey, id string) ([]string, error) {

	var gameParticipants GameData

	url := fmt.Sprintf("https://euw1.api.riotgames.com/lol/spectator/v4/active-games/by-summoner/%s", id)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("accept", "application/json")
	req.Header.Set("X-Riot-Token", apiKey)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return []string{}, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		fmt.Println(res)
		return []string{}, fmt.Errorf("received non-200 response code: %d", res.StatusCode)
	}
	body, _ := io.ReadAll(res.Body)
	err = json.Unmarshal(body, &gameParticipants)
	if err != nil {
		return []string{}, err
	}
	var participants []string
	for _, p := range gameParticipants.Participants {
		participants = append(participants, p.Puuid)
	}
	return participants, nil
}

func liveGamePlayersStats(apiKey, name, hashtag string) ([]LeagueEntry, error) {
	var playersStats []LeagueEntry
	puuID, err := GetPuuID(name, hashtag, apiKey)
	if err != nil {
		return playersStats, err
	}
	sumInf, err := getSummonerInfoByPuuID(puuID, apiKey, "euw1")
	if err != nil {
		return playersStats, err
	}
	players, err := LiveGamePlayersPuuIDS(apiKey, sumInf.ID)
	if err != nil {
		return playersStats, err
	}
	for _, p := range players {
		ps, err := RankedStats(p, apiKey, "euw1")
		if err != nil {
			return playersStats, err
		}
		playersStats = append(playersStats, ps)
	}
	return playersStats, nil
}
func LiveGamePlayersStatsFormattedToString(apiKey, name, hashtag string) (string, error) {
	stats, err := liveGamePlayersStats(apiKey, name, hashtag)
	if err != nil {
		return "", err
	}
	var totalLPBlue int
	var totalLPRed int
	var totalBIGLPBlue int
	var totalBIGLPRed int
	result := fmt.Sprintf("Active game players for %s \nBlue team\n", name)
	for i, s := range stats {
		if i == 5 {
			result += "\nRed Team\n"
		}
		if i < 5 {
			wr := float64(s.Wins) / float64(s.Wins+s.Losses) * 100
			result += fmt.Sprintf("%s: %s %s: %dlp Wins:%d Losses:%d wr:%.2f%%\n", s.SummonerName, s.Tier, s.Rank, s.LeaguePoints, s.Wins, s.Losses, wr)
			totalLPBlue += s.LeaguePoints
			lpGainz := utils.RankToLP(s.Tier, s.Rank, s.LeaguePoints)
			totalBIGLPBlue += lpGainz
		} else {
			wr := float64(s.Wins) / float64(s.Wins+s.Losses) * 100
			result += fmt.Sprintf("%s: %s %s: %dlp Wins:%d Losses:%d wr:%.2f%%\n", s.SummonerName, s.Tier, s.Rank, s.LeaguePoints, s.Wins, s.Losses, wr)
			totalLPRed += s.LeaguePoints
			lpGainz := utils.RankToLP(s.Tier, s.Rank, s.LeaguePoints)
			totalBIGLPRed += lpGainz
		}
	}
	result += fmt.Sprintf("Blue: total  %dlp avg %dlp Red: %dlp, %dlp\n", totalLPBlue, totalLPBlue/5, totalLPRed, totalLPRed/5)
	result += fmt.Sprintf("Blue total pisslow lp:%d, red:%dlp, difference:%d", totalBIGLPBlue, totalBIGLPRed, totalBIGLPBlue-totalBIGLPRed)
	return result, nil
}

type withName struct {
	puuID string
	name  string
}

func LiveGamePlayersPuuIDSv5(apiKey, id string) ([]withName, error) {

	var gameParticipants GameData
	var participants []withName
	url := fmt.Sprintf("https://euw1.api.riotgames.com/lol/spectator/v5/active-games/by-summoner/%s", id)

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("accept", "application/json")
	req.Header.Set("X-Riot-Token", apiKey)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return participants, err
	}
	defer res.Body.Close()
	if res.StatusCode != http.StatusOK {
		return participants, fmt.Errorf("received non-200 response code: %d", res.StatusCode)
	}
	body, _ := io.ReadAll(res.Body)
	err = json.Unmarshal(body, &gameParticipants)
	if err != nil {
		return participants, err
	}
	for _, p := range gameParticipants.Participants {
		var newP withName
		newP.puuID = p.Puuid
		newP.name = p.RiotId
		participants = append(participants, newP)
	}
	return participants, nil
}

func liveGamePlayersStatsPuuidSkip(apiKey, puuID string) ([]LeagueEntry, error) {
	var playersStats []LeagueEntry
	players, err := LiveGamePlayersPuuIDSv5(apiKey, puuID)
	if err != nil {
		return playersStats, err
	}
	for _, p := range players {
		ps, err := RankedStats(p.puuID, apiKey, "euw1")
		if err != nil {
			return playersStats, err
		}
		ps.SummonerName = p.name
		playersStats = append(playersStats, ps)
	}
	return playersStats, nil
}
func LiveGamePlayersStatsPuuIDSkipToString(apiKey, puuID, name string) (string, error) {
	stats, err := liveGamePlayersStatsPuuidSkip(apiKey, puuID)
	if err != nil {
		fmt.Println(err)
		return "", err
	}
	var totalLPBlue int
	var totalLPRed int
	var totalBIGLPBlue int
	var totalBIGLPRed int
	result := fmt.Sprintf("%s is now in game! \nBlue team\n", name)
	for i, s := range stats {
		if i == 5 {
			result += "\nRed Team\n"
		}
		if i < 5 {
			wr := float64(s.Wins) / float64(s.Wins+s.Losses) * 100
			result += fmt.Sprintf("%s: %s %s: %dlp Wins:%d Losses:%d wr:%.2f%%\n", s.SummonerName, s.Tier, s.Rank, s.LeaguePoints, s.Wins, s.Losses, wr)
			totalLPBlue += s.LeaguePoints
			lpGainz := utils.RankToLP(s.Tier, s.Rank, s.LeaguePoints)
			totalBIGLPBlue += lpGainz
		} else {
			wr := float64(s.Wins) / float64(s.Wins+s.Losses) * 100
			result += fmt.Sprintf("%s: %s %s: %dlp Wins:%d Losses:%d wr:%.2f%%\n", s.SummonerName, s.Tier, s.Rank, s.LeaguePoints, s.Wins, s.Losses, wr)
			totalLPRed += s.LeaguePoints
			lpGainz := utils.RankToLP(s.Tier, s.Rank, s.LeaguePoints)
			totalBIGLPRed += lpGainz
		}
	}
	result += fmt.Sprintf("Master+ Blue: total  %dlp avg %dlp | Red: %dlp avg: %dlp\n", totalLPBlue, totalLPBlue/5, totalLPRed, totalLPRed/5)
	result += fmt.Sprintf("Below master lp: Blue total lp:%d, | red:%dlp, difference:%d", totalBIGLPBlue, totalBIGLPRed, totalBIGLPBlue-totalBIGLPRed)
	return result, nil
}
