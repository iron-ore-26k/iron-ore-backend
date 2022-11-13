module ore

go 1.15

require (
	github.com/iron-ore-26k/ore-pb-gen v1.0.0
	google.golang.org/grpc v1.50.1
	google.golang.org/grpc/cmd/protoc-gen-go-grpc v1.2.0 // indirect
	google.golang.org/protobuf v1.28.1 // indirect
)

replace github.com/iron-ore-26k/ore-pb-gen v1.0.0 => /home/jgme/iron_ore_backend/pb/gen
