package sign

type SignProvider interface{
	LoadSignRewardConfig() (*[]map[string][]byte, error)
	LoadSignInfo(gameId string, openId string, signType int) (*[]map[string][]byte, error)
}
