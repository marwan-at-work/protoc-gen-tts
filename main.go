package main

import (
	pgs "github.com/lyft/protoc-gen-star"
	"google.golang.org/protobuf/types/pluginpb"
	"marwan.io/protoc-gen-tts/tts"
)

func main() {
	x := uint64(pluginpb.CodeGeneratorResponse_FEATURE_PROTO3_OPTIONAL)
	pgs.Init(pgs.DebugEnv("DEBUG"), pgs.SupportedFeatures(&x)).
		RegisterModule(tts.New()).
		RegisterPostProcessor(tts.NewFormatter()).
		Render()
}
