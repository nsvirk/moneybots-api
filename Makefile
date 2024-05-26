MAIN_PACKAGE_PATH=/home/ec2-user/apps/moneybotsapi
BINARY_NAME=moneybotsapi
SERVICE_EXE_PATH=/usr/local/bin


fetch:  # Fetches the  git repository
        @echo "---------------------------------------------------"
        @echo "==> Fetching repository from github.com"
        @echo "---------------------------------------------------"
       	# @rm -rf ${MAIN_PACKAGE_PATH}
       	# @git clone https://github.com/nsvirk/moneybotsapi.git
        @cd ${MAIN_PACKAGE_PATH} && git pull

tidy:  ## Cleans the Go module.
        @echo "---------------------------------------------------"
        @echo "==> Tidying module"
        @echo "---------------------------------------------------"
        @cd  ${MAIN_PACKAGE_PATH} && go mod tidy

build: ## Builds the executable
        @echo "---------------------------------------------------"
        @echo "==> Building executable "
        @echo "---------------------------------------------------"
        @cd  ${MAIN_PACKAGE_PATH} && go build -o ${BINARY_NAME} main.go
        @echo "Done"
        @echo "---------------------------------------------------"

deploy: ## Deploy to the /usr/local/bin folder
        @echo "---------------------------------------------------"
        @echo "==> Deploying executable "
        @echo "---------------------------------------------------"
        @sudo systemctl stop ${BINARY_NAME}
        @sudo rm -f ${SERVICE_EXE_PATH}/${BINARY_NAME}
        @sudo mv ${MAIN_PACKAGE_PATH}/${BINARY_NAME} ${SERVICE_EXE_PATH}
        @sudo systemctl daemon-reload
        @sudo systemctl start ${BINARY_NAME}

all:  ## Does all the comands
        @make fetch
        @make tidy
        @make build
        @make deploy
        @echo "---------------------------------------------------"
        @echo "==> Done "
        @echo "---------------------------------------------------"
        @sudo systemctl status ${BINARY_NAME}
        @echo "---------------------------------------------------"

catlog: ## Read the logs
	@cat /home/ec2-user/logs/moneybotsapi/log.log

caterr: ## Read the err log
	@cat /home/ec2-user/logs/moneybotsapi/error.log