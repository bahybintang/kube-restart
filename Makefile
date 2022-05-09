run:
	CONFIG_PATH=./kube-restart.yml go run main.go --kubeconfig="/home/$(shell whoami)/.kube/config"

build-image:
	docker build -t bintangbahy/kube-restart .

build:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o service .
