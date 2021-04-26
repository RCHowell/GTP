PACKAGE_ROOT=.

gen:
	protoc \
		./gtp.proto -I $(PACKAGE_ROOT)\
		--go_out=$(PACKAGE_ROOT)\
		--go_opt=paths=source_relative\
		--go-grpc_out=$(PACKAGE_ROOT)\
		--go-grpc_opt=paths=source_relative

clean:
	rm *.pb.go

.PHONY: clean gen
