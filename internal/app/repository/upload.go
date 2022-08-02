package repository

import (
	"context"
	"github.com/minio/minio-go/v7"
)

func UploadFile(c *minio.Client, objectName, filePath string) (minio.UploadInfo, error) {
	ctx := context.Background()
	bucketName := "android"

	return c.FPutObject(
		ctx,
		bucketName,
		objectName,
		filePath,
		minio.PutObjectOptions{
			ContentType: "application/vnd.android.package-archive",
		})
}
