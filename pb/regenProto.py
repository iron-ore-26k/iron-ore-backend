#!/usr/bin/env python3

import subprocess

print("Regenerating Proto")
cmd = "protoc --go_out=pb/gen --go_opt=paths=source_relative \
    --go-grpc_out=pb/gen --go-grpc_opt=paths=source_relative \
    pb/protobuf/songs.proto"

cmd = r"protoc --go_out=.\gen --go_opt=paths=source_relative --go_grpc_out=.\gen --proto_path=.\protobuf\ore -I.\protobuf\ore songs.proto"
# cmd2 = r"protoc --go_grpc_out=.\pb\gen --go-grpc_opt=paths=source_relative --proto_path=.\pb\protobuf\ore songs.proto service.proto"

subprocess.run(cmd)
# subprocess.run(cmd2)

print("Finished")