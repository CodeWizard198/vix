package vix

var (
	useNumber             = false
	disallowUnknownFields = false
)

type jsonMod interface {
	BindJSON(target any) error
	BindJSONbyOpt(target any, numberUse bool, disallow bool) error
}

func UseNumber() {
	useNumber = true
}

func DisallowUnknownFields() {
	disallowUnknownFields = true
}
