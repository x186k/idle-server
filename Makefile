all: build
	docker push x186k/idle-server




build:
	docker build -t x186k/idle-server .




run: build
	docker run x186k/idle-server




test: build 
	docker run -v ${PWD}:/foo x186k/idle-server --input /foo/idle-media




serve: build
	docker run -p 8080:8080 x186k/idle-server

serve-curl: build
	curl -X POST --data-binary @idle-media http://localhost:8080/convert --output idle-clip.zip


