package gcp

import (
	"context"
	"fmt"
	"google.golang.org/api/run/v2"
	"log"
)

func CreateCloudRunSilo() {
	ctx := context.Background()
	svc, err := run.NewService(ctx)
	if err != nil {
		log.Panic(err)
	}

	project := "projects/marcell-424212"
	location := "locations/asia-southeast2"
	service_name := "sample-saas-product"

	fullServiceName := fmt.Sprintf("%s/%s", project, location)
	// 11
	serviceObj := &run.GoogleCloudRunV2Service{
		// 11 12 lah sama terraform
		Ingress: "INGRESS_TRAFFIC_ALL",
		Template: &run.GoogleCloudRunV2RevisionTemplate{
			Containers: []*run.GoogleCloudRunV2Container{
				{
					Image: "asia-southeast2-docker.pkg.dev/marcell-424212/sample-saas-product/sample-saas-product:latest",
					Env: []*run.GoogleCloudRunV2EnvVar{
						{
							Name:  "DB_HOST",
							Value: "34.101.221.46",
						},
						{
							Name:  "DB_PORT",
							Value: "5432",
						},
						{
							Name:  "DB_DRIVER",
							Value: "postgres",
						},
						{
							Name:  "DB_DATABASE",
							Value: "postgres",
						},
						{
							Name:  "DB_USER",
							Value: "tenant-marsel",
						},
						{
							Name:  "DB_PASSWORD",
							Value: "Iron12345",
						},
					},
				},
			},
		},
	}
	_, err = svc.Projects.Locations.Services.Create(fullServiceName, serviceObj).ServiceId(service_name).Context(ctx).Do()
	if err != nil {
		log.Fatalf("Failed to deploy service: %v", err)
	}

	fmt.Printf("Cloud Service deployed")

	// setting IAM
	fullServiceName = fullServiceName + "/services/" + service_name
	policy, err := svc.Projects.Locations.Services.GetIamPolicy(fullServiceName).Context(ctx).Do()
	if err != nil {
		log.Fatalf("Failed to get IAM policy: %v", err)
	}

	binding := &run.GoogleIamV1Binding{
		Role:    "roles/run.invoker",
		Members: []string{"allUsers"},
	}
	policy.Bindings = append(policy.Bindings, binding)

	// Set the updated IAM policy
	setPolicyReq := &run.GoogleIamV1SetIamPolicyRequest{
		Policy: policy,
	}

	_, err = svc.Projects.Locations.Services.SetIamPolicy(fullServiceName, setPolicyReq).Context(ctx).Do()
	if err != nil {
		return
	}
	fmt.Printf("Cloud Iam Service deployed")
}
