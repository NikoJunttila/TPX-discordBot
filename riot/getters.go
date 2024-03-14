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

func GetMatch(matchId string, puuID string, country string, apiKey string) (string, bool, error) {
	var matchInfo MatchData
	endPoint := fmt.Sprintf("https://%s.api.riotgames.com/lol/match/v5/matches/%s", country, matchId)
	matchReq2, _ := http.NewRequest("GET", endPoint, nil)
	matchReq2.Header.Set("X-Riot-Token", apiKey)
	client := &http.Client{}
	matchRes2, err := client.Do(matchReq2)
	if err != nil {
		fmt.Println("Error making HTTP request:", err)
		return "", false, err
	}
	defer matchRes2.Body.Close()
	if matchRes2.StatusCode != http.StatusOK {
		fmt.Printf("Error get match: %s\n", matchRes2.Status)
		return "", false, err
	}
	body2, err := io.ReadAll(matchRes2.Body)
	if err != nil {
		fmt.Println("Error copying response body:", err)
		return "", false, err
	}
	// Unmarshal the JSON data into the struct
	err = json.Unmarshal(body2, &matchInfo)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return "", false, err
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
			result = fmt.Sprintf("New %s game %dmin! \nPlayer: %s\nRole: %s , Champion: %s, lvl: %d\nK/D/A: %d/%d/%d\nVisionWards placed: %d. Wards total placed: %d\n", queType, matchInfo.Info.GameDuration/60, p.RiotName, p.RoleNew, p.ChampionName, p.ChampLevel, p.Kills, p.Deaths, p.Assists, p.Wards, p.WardsPlaces)
			//result += fmt.Sprintf("Ally jg camps stolen: %d, Enemy camps stolen: %d\n",p.JgCampsStolen, p.EnemyJGCampsStolen,)
			/* 			if p.TimeSpentDead > 60 {
			   				result += fmt.Sprintf("Time wasted on dying %0.2fmin.\n", float32(p.TimeSpentDead)/60)
			   			} else {
			   				result += fmt.Sprintf("Time spent with gray screen %ds.\n", p.TimeSpentDead)
			   			} */
			result += fmt.Sprintf("Pings OMW: %d, KYS: %d, Missing: %d, GetBack: %d Danger: %d.\n", p.OnMyWayPings, p.KysPing, p.MissingPing, p.GetBackPings, p.DangerPing)
			result += fmt.Sprintf("Total dmg: %d , per min: %.2f, \n", p.DmgDealt, p.Challenges.DmgPerMinute)
			result += fmt.Sprintf("Lane minions first 10min: %d \n", p.Challenges.MinionsFirst10)
			result += fmt.Sprintf("Solo BOLOS: %d \n", p.Challenges.SoloBolo)
			if p.Win {
				result += "GG WP won game\n\n"
			} else {
				result += fmt.Sprintf("Game lost: GG %s GAP GIT GUUD NOOB\n\n", p.RoleNew)
			}

			fmt.Println(result)
		}
	}
	ranked := false
	if que == 420 {
		ranked = true
	}
	return result, ranked, nil
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
		fmt.Printf("Error getPuuId: %s\n", matchRes2.Status)
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

type CountingStats struct {
	ownCamps    int
	enemyCamps  int
	win         int
	deadTime    int
	deaths      int
	kills       int
	wardsPlaced int
	wardsBought int
	pings       int
	assists     int
	soloKills   int
	laneMinions int
	//laneAdvantage int
	dmgPermin float32
}

func getMatchStats(matchId string, puuID string, country string, apiKey string) (CountingStats, error) {
	var c CountingStats
	var matchInfo MatchData
	endPoint := fmt.Sprintf("https://%s.api.riotgames.com/lol/match/v5/matches/%s", country, matchId)
	matchReq2, _ := http.NewRequest("GET", endPoint, nil)
	matchReq2.Header.Set("X-Riot-Token", apiKey)
	client := &http.Client{}
	matchRes2, err := client.Do(matchReq2)
	if err != nil {
		fmt.Println("Error making HTTP request:", err)
		return c, err
	}
	defer matchRes2.Body.Close()
	if matchRes2.StatusCode != http.StatusOK {
		fmt.Printf("Error get match stats: %s\n", matchRes2.Status)
		return c, err
	}
	body2, err := io.ReadAll(matchRes2.Body)
	if err != nil {
		fmt.Println("Error copying response body:", err)
		return c, err
	}
	// Unmarshal the JSON data into the struct
	err = json.Unmarshal(body2, &matchInfo)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return c, err
	}
	for _, p := range matchInfo.Info.Participants {
		if p.Puuid == puuID {
			c.ownCamps = p.JgCampsStolen
			c.enemyCamps = p.EnemyJGCampsStolen
			c.win = 0
			if p.Win {
				c.win = 1
			}
			c.deadTime = p.TimeSpentDead
			c.deaths = p.Deaths
			c.kills = p.Kills
			c.assists = p.Assists
			c.wardsBought = p.Wards
			c.wardsPlaced = p.WardsPlaces
			c.pings = p.GetBackPings + p.OnMyWayPings + p.KysPing + p.MissingPing + p.DangerPing
			c.soloKills = p.Challenges.SoloBolo
			c.laneMinions = p.Challenges.MinionsFirst10
			c.dmgPermin = p.Challenges.DmgPerMinute
		}
	}
	return c, nil
}
func PrintHistory(matchHistory []string, apiKey string, puuID string, country string, playerName string) (string, error) {
	var c CountingStats
	for i, g := range matchHistory {
		inf, err := getMatchStats(g, puuID, country, apiKey)
		if err != nil {
			fmt.Println(err)
			break
		}
		fmt.Println(i)
		c.ownCamps += inf.ownCamps
		c.enemyCamps += inf.enemyCamps
		c.win += inf.win
		c.deadTime += inf.deadTime
		c.deaths += inf.deaths
		c.kills += inf.kills
		c.assists += inf.assists
		c.wardsBought += inf.wardsBought
		c.wardsPlaced += inf.wardsPlaced
		c.pings += inf.pings
		c.soloKills += inf.soloKills
		c.laneMinions += inf.laneMinions
		c.dmgPermin += inf.dmgPermin
	}
	var result string
	result = fmt.Sprintf("Avg stats for %s in last %d games \n", playerName, len(matchHistory))
	result += fmt.Sprintf("own team camps stolen: %d, enemy team: %d \nWins: %d, deadge timer: %0.2f minutes \n", c.ownCamps, c.enemyCamps, c.win, float32(c.deadTime)/60)
	result += fmt.Sprintf("Kills %d Deaths %d Assists %d and avg %.2f/%.2f/%.2f \n", c.kills, c.deaths, c.assists, float32(c.kills)/float32(len(matchHistory)), float32(c.deaths)/float32(len(matchHistory)), float32(c.assists)/float32(len(matchHistory)))
	result += fmt.Sprintf("Wards bought: %d, placed: %d \n", c.wardsBought, c.wardsPlaced)
	result += fmt.Sprintf("Total times pinged %d and avg per game %.2f \n", c.pings, float32(c.pings)/float32(len(matchHistory)))
	result += fmt.Sprintf("Solo Bolos: %d, DMG per min %.2f, minions first 10min avg: %.2f \n", c.soloKills, c.dmgPermin/float32(len(matchHistory)), float32(c.laneMinions)/float32(len(matchHistory)))
	result += fmt.Sprintf("WR = %.2f%%\n", float32(c.win)/float32(len(matchHistory))*100)

	return result, nil
}
