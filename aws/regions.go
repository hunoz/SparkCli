package aws

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
)

type Endpoints struct {
	Partitions []Partition `json:"partitions"`
}

type Partition struct {
	PartitionName string                 `json:"partitionName"`
	Regions       map[string]interface{} `json:"regions"`
	Services      map[string]Service     `json:"services"`
}

type Service struct {
	Endpoints map[string]interface{} `json:"endpoints"`
}

func GetAwsRegions() ([]string, error) {
	resp, err := http.Get("https://raw.githubusercontent.com/aws/aws-sdk-go-v2/master/codegen/smithy-aws-go-codegen/src/main/resources/software/amazon/smithy/aws/go/codegen/endpoints.json")
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
	}

	endpoints := Endpoints{}
	err = json.Unmarshal(body, &endpoints)
	if err != nil {
		return nil, err
	}

	var regions []string

	for _, parition := range endpoints.Partitions {
		for k := range parition.Regions {
			regions = append(regions, k)
		}
	}
	return regions, nil
}
