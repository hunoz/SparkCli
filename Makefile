make GOOS = "${GOOS}"
make GOARCH = "${GOARCH}"

PKG = gtech.dev/spark

ifndef GOOS
GOOS = darwin
endif

ifndef GOARCH
GOARCH = arm64
endif


build:
	@GOOS=${GOOS} GOARCH=${GOARCH} go build -o spark-${GOOS}-${GOARCH}

run:
	@GOOS=${GOOS} GOARCH=${GOARCH} go run main.go $(ARGS)