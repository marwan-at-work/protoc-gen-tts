package tts

type importData struct {
	Name         string
	Declarations []Declaration
}

type Declaration struct {
	IsType bool
	Name   string
}
