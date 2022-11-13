#!/usr/bin/env python3

import os
import subprocess

os.path.dirname(os.path.realpath(__file__)) 

print("Regenerating Proto")
cmd = "protoc --go_out=./gen --go_opt=paths=source_relative \
    --go-grpc_out=./gen --go-grpc_opt=paths=source_relative \
    protobuf/songs.proto"

cmd = r"protoc --go_out=.\gen --go_opt=paths=source_relative --go-grpc_out=.\gen --go-grpc_opt=paths=source_relative --proto_path=.\protobuf\ore -I.\ songs.proto service.proto"
if os.name != 'nt': # if not on windows
    cmd = cmd.replace('\\','/')
# cmd2 = r"protoc --go_grpc_out=.\pb\gen --go-grpc_opt=paths=source_relative --proto_path=.\pb\protobuf\ore songs.proto service.proto"

print(os.getcwd())
os.chdir(os.getcwd())
subprocess.Popen(cmd)
# subprocess.run(cmd2)

print("Finished")