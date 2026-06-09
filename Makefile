default: build push

build: 
	docker build -t ghcr.io/kuberhealthy/ssh-check:v3.0.0 .

push: 
	docker push ghcr.io/kuberhealthy/ssh-check:v3.0.0
