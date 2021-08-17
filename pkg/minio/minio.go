package main

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/go-logr/logr"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

// minioStorage
type minioStorage struct {
	// Context carries values across API boundaries.
	context.Context
	// Logger represents the ability to log messages.
	logr.Logger

	minioClient *minio.Client
}

func fileupload() error {
	minioClient, err := minio.New("10.108.88.174:9000", &minio.Options{
		Creds: credentials.NewStaticV4("AKIAIOSFODNN7EXAMPLE", "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY", ""),
		// Secure: true,
	})
	if err != nil {
		return err
	}
	// Make a new bucket.
	ctx := context.Background()
	bucketName := "mymusic"

	// err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{})
	// if err != nil {
	// 	// Check to see if we already own this bucket (which happens if you run this twice)
	// 	exists, errBucketExists := minioClient.BucketExists(ctx, bucketName)
	// 	if errBucketExists == nil && exists {
	// 		fmt.Println("bucket already exist.", "bucketName", bucketName)
	// 	} else {
	// 		return err
	// 	}
	// } else {
	// 	fmt.Println("successfully created bucket.", "bucketName", bucketName)
	// }

	// Upload the zip file
	objectName2 := "go-profiler-2938jddl23.my-node.0.0.0.20210728084437.default.go-profiler.Heap.prof"
	filePath2 := "/var/lib/kubediag/profilers/go/pprof/20210728084437/default.go-profiler.Heap.prof"
	contentType := "application/zip"
	// Upload the zip file with FPutObject
	info, err := minioClient.FPutObject(ctx, bucketName, objectName2, filePath2, minio.PutObjectOptions{ContentType: contentType})
	if err != nil {
		fmt.Println(err)
		return err
	}
	fmt.Println("Successfully uploaded object", "objectName", objectName2)
	fmt.Println("Successfully uploaded object", "infoKey", info)
	// AAA := strings.NewReader("AAAA")

	// minioClient.PutObject(ctx, bucketName, "go-profiler-2938jdl23.my-node.0.0.0.202007271149.key", AAA, AAA.Size(), minio.PutObjectOptions{ContentType: contentType})

	err = minioClient.FGetObject(context.Background(), bucketName, objectName2, "/home/shujiang/go/src/helloworld/profiler.pprof", minio.GetObjectOptions{})
	fmt.Println(err)

	// Set request parameters for content-disposition.
	reqParams := make(url.Values)
	reqParams.Set("response-content-disposition", "attachment; filename=\"your-filename.txt\"")

	// Generates a presigned url which expires in 7 days.
	presignedURL, err := minioClient.PresignedGetObject(ctx, bucketName, objectName2, time.Second*24*60*60*7, reqParams)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("Successfully generated presigned URL", presignedURL)

	// inf := minioClient.ListObjects(context.Background(), "mymusic", minio.ListObjectsOptions{})

	// a := <-inf
	// fmt.Println(a.Key)

	// err = minioClient.RemoveObject(context.Background(), "mymusic", "defaultpprof", minio.RemoveObjectOptions{})
	// if err != nil {
	// 	fmt.Println(err)
	// } else {
	// 	fmt.Println("object is deleted")
	// }
	// err = minioClient.RemoveBucket(context.Background(), "mymusic")
	// if err != nil {
	// 	fmt.Println(err)
	// } else {
	// 	fmt.Println("mymusic deleted")
	// }

	// infos, err = minioClient.ListBuckets(context.Background())
	// for _, ii := range infos {
	// 	fmt.Println(ii.Name)
	// }

	return nil
}
