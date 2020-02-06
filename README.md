### protoc-gen-tts

A Protobuf plugin that generates a Typescript client for Twirp services

This repository is based on https://github.com/horizon-games/protoc-gen-twirp_ts and https://github.com/larrymyers/protoc-gen-twirp_typescript (MIT)

It differs in that it correctly handles json numbers as strings and it includes documentation from protobuf comments into the generated TS file.

Current unsupported features (that I don't yet need) are nested messages, and maps. 

Status: works for my purposes but feel free to open issues or contribute.
