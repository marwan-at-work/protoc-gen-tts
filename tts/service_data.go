package tts

import pgs "github.com/lyft/protoc-gen-star"

type serviceData struct {
	Name       string
	Doc        string
	PathPrefix string
	Methods    []*methodData
}

type methodData struct {
	Name           string
	Doc            string
	Input          string
	Output         string
	InputHasFields bool
}

func createMethod(m pgs.Method) *methodData {
	return &methodData{
		Name:           m.Name().LowerCamelCase().String(),
		Doc:            getDoc(m.SourceCodeInfo().LeadingComments(), 1),
		Input:          m.Input().Name().String(),
		Output:         m.Output().Name().String(),
		InputHasFields: len(m.Input().Fields()) > 0,
	}
}
