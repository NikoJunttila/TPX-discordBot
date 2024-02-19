package riot

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

func InitStats(name string, tagLine string, country string, apiKey string) (string, string) {

	puuID, err := GetPuuID(name, tagLine, apiKey)
	if err != nil {
		fmt.Println("error getting puuid,", err)
		return "", ""
	}
	history, err := GetMatchHistory(puuID, 1, country, apiKey)
	if err != nil {
		fmt.Println(err, "err getting match history")
		return "", ""
	}
	return history[0], puuID
}

func GetMatch(matchId string, puuID string, country string, apiKey string) (string, error) {
	var matchInfo MatchData
	endPoint := fmt.Sprintf("https://%s.api.riotgames.com/lol/match/v5/matches/%s", country, matchId)
	matchReq2, _ := http.NewRequest("GET", endPoint, nil)
	matchReq2.Header.Set("X-Riot-Token", apiKey)
	client := &http.Client{}
	matchRes2, err := client.Do(matchReq2)
	if err != nil {
		fmt.Println("Error making HTTP request:", err)
		return "", err
	}
	defer matchRes2.Body.Close()
	if matchRes2.StatusCode != http.StatusOK {
		fmt.Printf("Error: %s\n", matchRes2.Status)
		return "", err
	}
	body2, err := io.ReadAll(matchRes2.Body)
	if err != nil {
		fmt.Println("Error copying response body:", err)
		return "", err
	}
	// Unmarshal the JSON data into the struct
	err = json.Unmarshal(body2, &matchInfo)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return "", err
	}
	var result string
	que := matchInfo.Info.QueueID
	var queType string
	switch que {
	case 400:
		queType = "Normal"
	case 420:
		queType = "Ranked Soloq"
	case 440:
		queType = "Ranked Flex"
	case 450:
		queType = "ARAM"
	default:
		queType = "UNKNOWN"
	}
	for _, p := range matchInfo.Info.Participants {
		if p.Puuid == puuID {
			result = fmt.Sprintf("New %s game! \nPlayer: %s\nRole: %s , Champion: %s, lvl: %d\nKills: %d, Deaths: %d Assists: %d\nAlly jg camps stolen: %d. Enemy camps stolen: %d\nWards bought: %d. Wards placed: %d\n", queType, p.RiotName, p.RoleNew, p.ChampionName, p.ChampLevel, p.Kills, p.Deaths, p.Assists, p.JgCampsStolen, p.EnemyJGCampsStolen, p.Wards, p.WardsPlaces)
			if p.TimeSpentDead > 60 {
				result += fmt.Sprintf("Time wasted on dying %0.2fmin.\n", float32(p.TimeSpentDead)/60)
			} else {
				result += fmt.Sprintf("Time spent with gray screen %ds.\n", p.TimeSpentDead)
			}
			result += fmt.Sprintf("Pings OMW: %d, KYS: %d, Missing: %d, GetBack: %d. \n", p.OnMyWayPings, p.KysPing, p.MissingPing, p.GetBackPings)
			if p.Win {
				result += "Boosted monkey won\n\n"
			} else {
				result += fmt.Sprintf("Game lost: GG %s GAP GIT GUUD NOOB\n\n", p.RoleNew)
			}

			fmt.Println(result)
		}
	}
	return result, nil
}

func GetPuuID(name string, tagLine string, apiKey string) (string, error) {
	type puuIDfromRiotID struct {
		PuuID string `json:"puuid"`
	}
	var puuId puuIDfromRiotID
	endPoint := fmt.Sprintf("https://europe.api.riotgames.com/riot/account/v1/accounts/by-riot-id/%s/%s", name, tagLine)
	matchReq2, _ := http.NewRequest("GET", endPoint, nil)
	matchReq2.Header.Set("X-Riot-Token", apiKey)
	client := &http.Client{}
	matchRes2, err := client.Do(matchReq2)
	if err != nil {
		fmt.Println("Error making HTTP request:", err)
		return "", err
	}
	defer matchRes2.Body.Close()
	if matchRes2.StatusCode != http.StatusOK {
		fmt.Printf("Error: %s\n", matchRes2.Status)
		return "", err

	}
	body2, err := io.ReadAll(matchRes2.Body)
	if err != nil {
		fmt.Println("Error copying response body:", err)
		return "", err
	}
	// Unmarshal the JSON data into the struct
	err = json.Unmarshal(body2, &puuId)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return "", err

	}
	return puuId.PuuID, nil
}
