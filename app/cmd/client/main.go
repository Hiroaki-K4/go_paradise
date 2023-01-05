import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	hellopb "mygrpc/pkg/grpc"
)

var (
	scanner *bufio.Scanner
	client hellopb.GreetingServiceClient
)

func main() {
	fmt.Println("Start gRPC client.")

	scanner = bufio.NewScanner(os.Stdin)

	address := "localhost:8080"
	conn, err := gcpc
}
