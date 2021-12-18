.PHONY: all dev clean build env-up env-down run dbinit

all: clean build env-up run

dev: build run

##### BUILD
build:
	@echo "Build ..."
	@go mod vendor
	@cd chaincode && go mod vendor
	@go build
	@echo "Build done"

##### ENV
env-up:
	@echo "Start environment ..."
	@cd ${GOPATH}/src/fabric-samples/test-network && ./network.sh up
	@echo "Environment up"

env-down:
	@echo "Stop environment ..."
	@cd ${GOPATH}/src/fabric-samples/test-network && ./network.sh down
	@echo "Environment down"

##### RUN
run:
	@echo "Start app ..."
	@cd ${GOPATH}/src/fabric-samples/test-network && configtxgen -profile TwoOrgsChannel -outputCreateChannelTx ./channel-artifacts/testchannel.tx -channelID testchannel  -configPath ./configtx/
	@cd ${GOPATH}/src/fabric-samples/test-network && configtxgen -profile TwoOrgsChannel -outputAnchorPeersUpdate ./channel-artifacts/Org1MSPanchors.tx -channelID testchannel -asOrg Org1MSP -configPath ./configtx/
	@cd ${GOPATH}/src/fabric-samples/test-network && configtxgen -profile TwoOrgsChannel -outputAnchorPeersUpdate ./channel-artifacts/Org2MSPanchors.tx -channelID testchannel -asOrg Org2MSP -configPath ./configtx/
	@chmod +x ./fabric-first-go-app
	@chmod +x ../fabric-samples/test-network/addOrg3/runme.sh
	@./fabric-first-go-app

##### CLEAN
clean: env-down
	@echo "Clean up ..."
	@rm -rf /home/verayy/data/fabric-first-go-app/*
	@cd ${GOPATH}/src/fabric-samples/test-network && ./network.sh down
	@echo "Clean up done ..."

### DBINIT 
dbinit:
	@rm -rf ./dbinit
	@cd ${GOPATH}/src/fabric-first-go-app/db/dbinit && go build && cp ./dbinit ../../
	@chmod +x ./dbinit
	@./dbinit

### WEBUP
webup:
	@cd ${GOPATH}/src/fabric-first-go-app/Webstart && go build && cp ./Webstart ../webstart
	@chmod +x ./webstart
	@./webstart


### airflow
