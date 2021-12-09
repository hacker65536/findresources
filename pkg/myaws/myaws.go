package myaws

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	rgtapi "github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi"
	"github.com/aws/aws-sdk-go-v2/service/resourcegroupstaggingapi/types"
	"github.com/iancoleman/strcase"
	log "github.com/sirupsen/logrus"
)

func ListAwsRes(tag map[string]string, restype []string) {
	// Load the Shared AWS Configuration (~/.aws/config)
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	t := []types.TagFilter{}
	for k, v := range tag {
		t = append(t, types.TagFilter{
			Key:    aws.String(strcase.ToCamel(k)),
			Values: []string{v},
		},
		)
	}

	/*
		restype := []string{

			// "ec2:instance",
			// "rds:db",
			"elasticloadbalancing:loadbalancer",
		}
		tag := []types.TagFilter{
			{
				Key:    aws.String("Stage"),
				Values: []string{"production"},
			},
			{
				Key:    aws.String("Project"),
				Values: []string{"myproject"},
			},
		}
	*/

	svc := rgtapi.NewFromConfig(cfg)

	params := &rgtapi.GetResourcesInput{
		ResourceTypeFilters: restype,
		TagFilters:          t,
	}

	log.WithFields(
		log.Fields{
			"params": params,
		}).Debug("[debug]")
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
