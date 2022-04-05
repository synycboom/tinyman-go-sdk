SRC_PATH := $(shell pwd)

generate:
	cd $(SRC_PATH) && go generate ./v1/contracts
