# docker build -t bme280
# docker run --device /dev/i2c-1 -d --rm -p 2112:2112 bme280
FROM golang:alpine

COPY . /src 
WORKDIR /src

ENV GO111MODULE=on
RUN CGO_ENABLED=0 go build -o bme280 

CMD ["./bme280"]

FROM scratch
COPY --from=0 /src/bme280 /bme280
CMD ["/bme280"]