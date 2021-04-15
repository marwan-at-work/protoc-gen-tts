# protoc-gen-tts

A Protobuf plugin that generates a Typescript client for Twirp services.

This protoc plugin is intended to work with the official JS protobuf plugin and
our related <https://github.com/InVisionApp/twirpjs>. It will generate js
service clients by wrapping `twirpjs` and referencing the generated protobuf
code from the official JS plugin when referencing messages. It also generates
type annotations for the generated JS code.

## Known Limitations

- Method inputs and outputs must be defined as messages in the protobuf package
  being rendered. For example,
  `rpc Do(google.protobuf.Empty) returns (google.protobuf.Empty)` is not yet
  handled. The inputs and outputs can still make use of imported types but the
  in/out types, themselves, must be framed in the local package like
  `rpc Do(DoRequest) returns (DoResponse)`.

