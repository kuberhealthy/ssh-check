default: build push

build: 
	docker build -t rjacks161/ssh-check:v1.0.0 .

push: 
	docker push rjacks161/ssh-check:v1.0.0 
