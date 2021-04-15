package tts

import _ "embed"

//go:embed template.txt
var jsTemplate string

//go:embed template.d.txt
var tsTemplate string
