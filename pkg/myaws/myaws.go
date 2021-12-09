package myaws

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	rgtapi "github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi"
	"github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi/types"
	log "github.com/sirupsen/logrus"
)

type TagFilter struct {
	Key string
	Val string
}

func ListAwsRes(tag []TagFilter, restypes []string) {
	// Load the Shared AWS Configuration (~/.aws/config)
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	tfs := []types.TagFilter{}
	for _, v := range tag {
		tfs = append(tfs, types.TagFilter{
			Key:    aws.String(v.Key),
			Values: []string{v.Val},
		})
	}

	svc := rgtapi.NewFromConfig(cfg)

	params := &rgtapi.GetResourcesInput{
		ResourceTypeFilters: restypes,
		TagFilters:          tfs,
	}

	log.WithFields(
		log.Fields{
			"params": func() string {
				j, _ := json.Marshal(params)
				return string(j)
			}(),
		}).Debug()
	p := rgtapi.NewGetResourcesPaginator(svc, params, func(o *rgtapi.GetResourcesPaginatorOptions) {
		o.Limit = 100
		o.StopOnDuplicateToken = true
	})

	pageNum := 0
	for p.HasMorePages() && pageNum < 10 {

		resp, err := p.NextPage(context.TODO())
		if err != nil {
			log.Printf("error: %v", err)
			return
		}
		for _, v := range resp.ResourceTagMappingList {
			fmt.Println(aws.ToString(v.ResourceARN))
		}
		// bug https://github.com/aws/aws-sdk-go-v2/issues/1201
		if *resp.PaginationToken == "" {
			break
		}
		pageNum++
	}
}
