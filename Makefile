.PHONY:docker delete deploy
docker:
	go env -w GOOS=linux
	go env -w GOARCH=amd64
	go build -o webook .
	docker rmi yakumo/webook:0.0.1
	docker build -t yakumo/webook:0.0.1 .
	go env -w GOOS=windows
delete:
	kubectl delete service webook
	kubectl delete deployment webook
deploy:
	kubectl apply -f .\k8s-webook-deployment.yaml
	kubectl apply -f .\k8s-webook-service.yaml
