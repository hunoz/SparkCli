make SPARK_CLIENT_ID = "${SPARK_CLIENT_ID}"
make SPARK_POOL_ID = "${SPARK_POOL_ID}"
make SPARK_POOL_REGION = "${SPARK_POOL_REGION}"
make GOOS = "${GOOS}"
make GOARCH = "${GOARCH}"

PKG = "gtech.dev/spark"

ifndef SPARK_CLIENT_ID
$(error SPARK_CLIENT_ID is not set)
endif

ifndef SPARK_POOL_ID
$(error SPARK_POOL_ID is not set)
endif

ifndef SPARK_POOL_REGION
$(error SPARK_POOL_REGION is not set)
endif

ifndef GOOS
GOOS = "darwin"
endif

ifndef GOARCH
GOARCH = "arm"
endif

ClientIdLdFlag = "${PKG}/cognito.ClientId=${SPARK_CLIENT_ID}"
PoolIdLdFlag = "${PKG}/cognito.PoolId=${SPARK_POOL_ID}"
PoolRegionFlag = "${PKG}/cognito.PoolRegion=${SPARK_POOL_REGION}"


build:
	@GOOS=${GOOS} GOARCH=${GOARCH} go build -ldflags="-X ${ClientIdLdFlag} -X ${PoolIdLdFlag} ${PoolRegionFlag} -o spark-${GOOS}-${GOARCH}"

run:
	@GOOS=${GOOS} GOARCH=${GOARCH} go run main.go -ldflags="-X ${ClientIdLdFlag} -X ${PoolIdLdFlag} ${PoolRegionFlag}"