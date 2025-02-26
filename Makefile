gen:
	protoc --descriptor_set_out=gen/descriptor.pb --go_out=gen/ proto/*.proto

.PHONY:
	gen