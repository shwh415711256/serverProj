package game

type GameProvider interface{
	LoadOneMatchHis(gameid string, openid string) (*[]map[string][]byte, error)
	InsertOneMatchHis(gameid string, openid string, data string) error
}
