gen:
	protoc --descriptor_set_out=gen/descriptor.pb --go_out=gen/ proto/*.proto
genWithDoc:
	protoc --descriptor_set_out=gen/descriptor.pb --doc_out=. --doc_opt=json,documentation.json --go_out=gen/ proto/*.proto
genDecodeDescriptor:
	protoc --decode_raw < gen/descriptor.pb
.PHONY:
	gen