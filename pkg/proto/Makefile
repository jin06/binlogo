pipeline.pb.go: pipeline.proto
	protoc pipeline.proto --go_out=plugins=grpc:. --go_opt=paths=source_relative

force:
	rm -f pipeline.pb.go
	make pipeline.pb.go

