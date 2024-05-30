package xingoudoc

type ApiDoc struct {
	Head    DocHead              `json:"head"`
	Params  map[string]DocParam  `json:"params"`
	Returns map[string]DocReturn `json:"returns"`
}

type DocCommon interface {
	Analyze(string)
	check()
}
