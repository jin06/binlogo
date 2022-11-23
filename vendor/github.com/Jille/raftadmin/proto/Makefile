raftadmin.pb.go: raftadmin.proto
	protoc raftadmin.proto --go_out=plugins=grpc:. --go_opt=paths=source_relative

force:
	rm -f raftadmin.pb.go
	make raftadmin.pb.go
