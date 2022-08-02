package db

import (
	"github.com/FTChinese/superyard/pkg/config"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"log"
)

func NewMinIOClient(c config.MinIOConfig) (*minio.Client, error) {
	return minio.New(c.Endpoint, &minio.Options{
		Creds:        credentials.NewStaticV2(c.AccessKeyID, c.SecretAccessKey, ""),
		Secure:       true,
		Transport:    nil,
		Region:       "",
		BucketLookup: 0,
		CustomMD5:    nil,
		CustomSHA256: nil,
	})
}

func MustMinIOClient(c config.MinIOConfig) *minio.Client {
	client, err := NewMinIOClient(c)
	if err != nil {
		log.Fatal(err)
	}

	return client
}
