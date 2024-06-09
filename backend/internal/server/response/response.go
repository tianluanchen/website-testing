package response

type status struct {
	OK  int
	Bad int
}

var Status = &status{
	OK:  0,
	Bad: 1,
}

type Universal struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

type UniversalWithData struct {
	Universal
	Data any `json:"data"`
}

type Error = Universal
type ErrorWithData = UniversalWithData
