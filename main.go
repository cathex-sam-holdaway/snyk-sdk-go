package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/cathex-sam-holdaway/snyk-sdk-go/snyk"
	"github.com/joho/godotenv"
	"github.com/kelseyhightower/envconfig"
)

type Snyk struct {
	Client       *snyk.Client
	SnykApiToken string `envconfig:"SNYK_API_TOKEN" required:"true"`
}

func main() {
	var s Snyk

	err := godotenv.Load()
	if err != nil {
		fmt.Printf("error loading .env file : %v", err)
	}

	if err := envconfig.Process("", &s); err != nil {
		fmt.Printf("error getting environment variables : %v", err)
	}

	snykClient := snyk.NewClient(s.SnykApiToken)

	targets, _, err := snykClient.Targets.List(context.Background(), "23ec072e-580e-4e16-bd64-856d3bad9aac")

	if err != nil {
		fmt.Printf("error getting targets : %v", err)
	}

	output, _ := json.Marshal(targets)

	fmt.Println(string(output))
}
