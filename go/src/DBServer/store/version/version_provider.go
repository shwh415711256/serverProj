package version

type VersionProvider interface{
	LoadVersionConfigData() (*[]map[string][]byte, error)
}