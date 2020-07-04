.PHONY: build run

build:
	docker build -t bme280 .

run:
	docker run --device /dev/i2c-1 -d --rm -p 2112:2112 bme280