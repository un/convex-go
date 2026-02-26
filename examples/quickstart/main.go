package main

import (
	"context"
	"fmt"
	"os"

	"github.com/get-convex/convex-go/convex"
)

func main() {
	deploymentURL := os.Getenv("CONVEX_DEPLOYMENT_URL")
	if deploymentURL == "" {
		panic("set CONVEX_DEPLOYMENT_URL")
	}

	client := convex.NewClientBuilder().WithDeploymentURL(deploymentURL).WithClientID("quickstart").Build()
	defer client.Close()

	result, err := client.Query(context.Background(), "messages:list", map[string]any{})
	if err != nil {
		panic(err)
	}
	value, err := result.Unwrap()
	if err != nil {
		panic(err)
	}
	fmt.Printf("query result: %#v\n", value.Raw())
}
