all:
	echo hello world
	docker build -t x186k/idle-server .
	docker run x186k/idle-server
	docker push x186k/idle-server