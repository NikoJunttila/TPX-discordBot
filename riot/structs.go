package riot

type GameData struct {
	GameId            int           `json:"gameId"`
	MapId             int           `json:"mapId"`
	GameMode          string        `json:"gameMode"`
	GameType          string        `json:"gameType"`
	GameQueueConfigId int           `json:"gameQueueConfigId"`
	Participants      []Participant `json:"participants"`
}
type Participant struct {
	Puuid         string `json:"puuid"`
	TeamId        int    `json:"teamId"`
	Spell1Id      int    `json:"spell1Id"`
	Spell2Id      int    `json:"spell2Id"`
	ChampionId    int    `json:"championId"`
	ProfileIconId int    `json:"profileIconId"`
	SummonerName  string `json:"summonerName"`
	RiotId        string `json:"riotId"`
	Bot           bool   `json:"bot"`
	SummonerId    string `json:"summonerId"`
}

type SummonerInfo struct {
	ID            string `json:"id"`
	AccountID     string `json:"accountId"`
	Puuid         string `json:"puuid"`
	ProfileIconID int    `json:"profileIconId"`
	RevisionDate  int64  `json:"revisionDate"`
	SummonerLevel int    `json:"summonerLevel"`
}
type LeagueEntry struct {
	LeagueID     string `json:"leagueId"`
	QueueType    string `json:"queueType"`
	Tier         string `json:"tier"`
	Rank         string `json:"rank"`
	SummonerID   string `json:"summonerId"`
	SummonerName string `json:"summonerName"`
	LeaguePoints int    `json:"leaguePoints"`
	Wins         int    `json:"wins"`
	Losses       int    `json:"losses"`
	Inactive     bool   `json:"inactive"`
	FreshBlood   bool   `json:"freshBlood"`
	HotStreak    bool   `json:"hotStreak"`
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
	ChampLevel         int               `json:"champLevel"`
	ChampionName       string            `json:"championName"`
	Role               string            `json:"role"`
	Wards              int               `json:"detectorWardsPlaced"`
	Puuid              string            `json:"puuid"`
	RiotName           string            `json:"riotIdGameName"`
	RoleNew            string            `json:"teamPosition"`
	JgCampsStolen      int               `json:"totalAllyJungleMinionsKilled"`
	EnemyJGCampsStolen int               `json:"totalEnemyJungleMinionsKilled"`
	TimeSpentDead      int               `json:"totalTimeSpentDead"`
	WardsPlaces        int               `json:"wardsPlaced"`
	Win                bool              `json:"win"`
	Deaths             int               `json:"deaths"`
	Assists            int               `json:"assists"`
	Kills              int               `json:"kills"`
	GetBackPings       int               `json:"getBackPings"`
	OnMyWayPings       int               `json:"onMyWayPings"`
	KysPing            int               `json:"pushPings"`
	MissingPing        int               `json:"enemyMissingPings"`
	DangerPing         int               `json:"dangerPings"`
	DmgDealt           int               `json:"totalDamageDealtToChampions"`
	TimePlayed         int               `json:"timePlayed"`
	Challenges         ChallengesDetails `json:"challenges"`
}
type ChallengesDetails struct {
	DmgPerMinute   float32 `json:"damagePerMinute"`
	MinionsFirst10 int     `json:"laneMinionsFirst10Minutes"`
	LaneAdvantage  int     `json:"laningPhaseGoldExpAdvantage"`
	SoloBolo       int     `json:"soloKills"`
}
