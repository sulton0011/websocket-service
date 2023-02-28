CURRENT_DIR=$(shell pwd)



swag-init:
	swag init -g api/api.go -o api/docs