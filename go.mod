module github.com/cyruslo/helper

go 1.16

require (
	//git.huoys.com/qp/luabridge v1.0.0
	github.com/BurntSushi/toml v0.3.1
	github.com/bilibili/kratos v0.3.3
	github.com/cyruslo/proto v0.0.0-20220209101156-0dee368f0d12
	github.com/golang/protobuf v1.5.2
	google.golang.org/grpc v1.43.0
)

replace (
	github.com/golang/protobuf => github.com/golang/protobuf v1.3.2
)