all: build
	docker push x186k/idle-server




build: main.go go.mod go.sum
	docker build -t x186k/idle-server .




run: build
	docker run x186k/idle-server




test: build 
	docker run -v ${PWD}:/foo x186k/idle-server /foo/idle-media.mov /foo/idle-clip.zip

test-deadsfu: build
	docker run -v ~/Documents/deadsfu/deadsfu-binaries:/foo x186k/idle-server /foo/idle-clip.mov /foo/idle-clip.zip


serve: build
	docker run -p 8080:8088 x186k/idle-server

serve-curl:
	curl -X POST --data-binary @idle-media http://localhost:8088/idle-clip --output idle-clip.zip



