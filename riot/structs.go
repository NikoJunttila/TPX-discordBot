package riot

type SummonerInfo struct {
	ID            string `json:"id"`
	AccountID     string `json:"accountId"`
	Puuid         string `json:"puuid"`
	Name          string `json:"name"`
	ProfileIconID int    `json:"profileIconId"`
	RevisionDate  int64  `json:"revisionDate"`
	SummonerLevel int    `json:"summonerLevel"`
}
type LeagueEntry struct {
	LeagueID     string `json:"leagueId"`
	SummonerID   string `json:"summonerId"`
	SummonerName string `json:"summonerName"`
	QueueType    string `json:"queueType"`
	Tier         string `json:"tier"`
	Rank         string `json:"rank"`
	LeaguePoints int    `json:"leaguePoints"`
	Wins         int    `json:"wins"`
	Losses       int    `json:"losses"`
	Inactive     bool   `json:"inactive"`
	FreshBlood   bool   `json:"freshBlood"`
	HotStreak    bool   `json:"hotStreak"`
	MiniSeries   struct {
		Target   int    `json:"target"`
		Wins     int    `json:"wins"`
		Losses   int    `json:"losses"`
		Progress string `json:"progress"`
	} `json:"miniSeries"`
}

type MatchData struct {
	MetaData MatchMetaData `json:"metaData"`
	Info     MatchInfo     `json:"info"`
}

type MatchMetaData struct {
	DataVersion  string   `json:"dataVersion"`
	MatchId      string   `json:"matchId"`
	Participants []string `json:"participants"`
}

type MatchInfo struct {
	EndOfGameResult    string               `json:"endOfGameResult"`
	GameCreation       int64                `json:"gameCreation"`
	GameDuration       int                  `json:"gameDuration"`
	GameEndTimestamp   int64                `json:"gameEndTimestamp"`
	GameId             int64                `json:"gameId"`
	GameMode           string               `json:"gameMode"`
	GameName           string               `json:"gameName"`
	GameStartTimestamp int64                `json:"gameStartTimestamp"`
	GameType           string               `json:"gameType"`
	GameVersion        string               `json:"gameVersion"`
	MapId              int                  `json:"mapId"`
	Participants       []ParticipantDetails `json:"participants"`
	QueueID            int                  `json:"queueId"`
	// Add other fields as needed
}

type ParticipantDetails struct {
	ChampLevel         int    `json:"champLevel"`
	ChampionName       string `json:"championName"`
	Role               string `json:"role"`
	Wards              int    `json:"sightWardsBoughtInGame"`
	Puuid              string `json:"puuid"`
	RiotName           string `json:"riotIdGameName"`
	RoleNew            string `json:"teamPosition"`
	JgCampsStolen      int    `json:"totalAllyJungleMinionsKilled"`
	EnemyJGCampsStolen int    `json:"totalEnemyJungleMinionsKilled"`
	TimeSpentDead      int    `json:"totalTimeSpentDead"`
	WardsPlaces        int    `json:"wardsPlaced"`
	Win                bool   `json:"win"`
	Deaths             int    `json:"deaths"`
	Assists            int    `json:"assists"`
	Kills              int    `json:"kills"`
	GetBackPings       int    `json:"getBackPings"`
	OnMyWayPings       int    `json:"onMyWayPings"`
	KysPing            int    `json:"pushPings"`
	MissingPing        int    `json:"enemyMissingPings"`
	DangerPing         int    `json:"dangerPings"`
}
