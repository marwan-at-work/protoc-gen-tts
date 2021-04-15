package tts

import (
	"text/template"

	pgs "github.com/lyft/protoc-gen-star"
)

// New tts plugin
func New() pgs.Module {
	return &tts{
		ModuleBase: &pgs.ModuleBase{},
	}
}

type tts struct {
	*pgs.ModuleBase
}

// Name is the identifier used to identify the module. This value is
// automatically attached to the BuildContext associated with the ModuleBase.
func (t *tts) Name() string { return "tts" }

func (t *tts) Execute(targets map[string]pgs.File, pkgs map[string]pgs.Package) []pgs.Artifact {
	for _, f := range targets {
		fname := f.Name().LowerSnakeCase().String()
		fname = fname[:len(fname)-6] // strip _proto from the value.
		t.AddGeneratorTemplateFile(fname+"_twirp_pb.js", template.Must(template.New("js").Parse(jsTemplate)), f)
		t.AddGeneratorTemplateFile(fname+"_twirp_pb.d.ts", template.Must(template.New("ts").Parse(tsTemplate)), f)
	}
	return t.Artifacts()
}
