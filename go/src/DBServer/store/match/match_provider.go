package match

type MatchProvider interface{
	LoadMatchConfigData()(*[]map[string][]byte, error)
}
