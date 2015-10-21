
build:
	go build -o gopocsag

clean:
	rm gopocsag
	rm -rf batches

test:
	go test ./...

docgen: 
	godoc -html ./internal/datatypes/ > doc/datatypes.html
	godoc -html ./internal/utils/ > doc/utils.html
	godoc -html ./internal/wav/ > doc/wav.html