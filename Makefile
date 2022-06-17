default: build push

build: 
	docker build -t hub.comcast.net/k8s-eng/ssh-check:v1.0.0 .

push: 
	docker push hub.comcast.net/k8s-eng/ssh-check:v1.0.0 
