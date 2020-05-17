build:
	cd cmd/vkwatcher && CC=x86_64-linux-musl-gcc CXX=x86_64-linux-musl-g++ CGO_ENABLED=1 GOOS=linux go build -a -ldflags '-linkmode external -extldflags -static' -o ./build
