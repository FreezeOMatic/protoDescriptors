gen:
	protoc --descriptor_set_out=gen/descriptor.pb --go_out=gen/ proto/*.proto
genWithDoc:
	protoc --descriptor_set_out=gen/descriptor.pb --doc_out=. --doc_opt=json,documentation.json --go_out=gen/ proto/*.proto
genDecodeDescriptor:
	protoc --decode_raw < gen/descriptor.pb
run:
	go run main.go -col=A:parent -col=B:enum -col=C:enumVName -col=D:enumVID -col=E:extInfo -col=F:extDesc -col=G:extAbb
.PHONY:
	gen