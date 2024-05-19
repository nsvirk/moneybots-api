BINARY_NAME=moneybotsapi

deploy:
	sudo systemctl stop moneybotsapi
	rm -f /usr/local/bin/${BINARY_NAME}
	cd ~/apps
	rm -rf moneybotsapi
	git clone https://github.com/nsvirk/moneybotsapi.git
	cd moneybotsapi/
	go mod tidy
	go build -o /usr/local/bin/${BINARY_NAME} main.go
	sudo systemctl start moneybotsapi
	sudo systemctl enable moneybotsapi
	sudo systemctl status moneybotsapi




