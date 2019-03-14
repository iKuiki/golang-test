if  [ ! $TARGET_PROJECT ];then
	TARGET_PROJECT="github.com/ikuiki/golang-test"
fi
echo "TARGET_PROJECT: $TARGET_PROJECT"

GRPC_PATHS=(
	"$TARGET_PROJECT/grpc-conn-create-benchmark/helloworld"
)

PBF_PATHS=(
)
