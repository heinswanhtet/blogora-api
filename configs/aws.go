package configs

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

var S3Client *s3.Client = setupS3Client()

func setupS3Client() *s3.Client {
	accessKeyID := Envs.AWS_ACCESS_KEY_ID
	secretAccessKey := Envs.AWS_SECRET_ACCESS_KEY
	region := Envs.AWS_REGION

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKeyID, secretAccessKey, "")),
	)
	if err != nil {
		panic(fmt.Sprintf("unable to load SDK config, %v", err))
	}

	return s3.NewFromConfig(cfg)
}
