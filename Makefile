all: pcp

pcp:
	go build -o pcp ./parallel_copy

run_tests: clean pcp
	go test -count=1 -v ./threadpool && \
	go test -count=1 -v ./fileutil 

clean:
	rm -f pcp
