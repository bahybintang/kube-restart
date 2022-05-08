run:
	CONFIG_PATH=./kube-restart.yml go run main.go --kubeconfig="/home/$(shell whoami)/.kube/config"

build:
	docker build -t bintangbahy/kube-restart .
