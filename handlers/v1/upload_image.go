package v1

import (
	"context"
	"fit-byte/config"
	"fit-byte/utils"
	"fmt"
	"log"
	"mime/multipart"
	"net/http"
	"sync"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/gin-gonic/gin"
)

type FileController struct {
	S3Client   *s3.Client
	BucketName string
}

func NewFileController(s3Client *s3.Client) *FileController {
	return &FileController{
		S3Client:   s3Client,
		BucketName: config.LoadConfig().S3Bucket,
	}
}

func InitS3Client() *s3.Client {
	cfg, err := awsConfig.LoadDefaultConfig(context.TODO(),
		awsConfig.WithRegion(config.LoadConfig().AwsRegion),
		awsConfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(config.LoadConfig().AwsAccessKeyId, config.LoadConfig().AwsSecretAccessKey, "")),
	)
	if err != nil {
		log.Fatalf("unable to load AWS SDK config, %v", err)
	}
	return s3.NewFromConfig(cfg)
}

func UploadToS3(s3Client *s3.Client, fileHeader *multipart.FileHeader, bucketName string, resultChan chan<- string) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	file, err := fileHeader.Open()
	if err != nil {
		resultChan <- ""
		return
	}
	defer file.Close()

	key := fmt.Sprintf("uploads/%d_%s", time.Now().UnixNano(), fileHeader.Filename)

	_, err = s3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(key),
		Body:        file,
		ACL:         types.ObjectCannedACLPublicRead,
		ContentType: aws.String(fileHeader.Header.Get("Content-Type")),
	})
	if err != nil {
		resultChan <- ""
		return
	}

	resultChan <- fmt.Sprintf("https://%s.s3.amazonaws.com/%s", bucketName, key)
}

func (fc *FileController) UploadFile(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File is required"})
		return
	}

	if err := utils.ValidateFile(file); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var wg sync.WaitGroup
	resultChan := make(chan string, 1)

	wg.Add(1)
	go func() {
		defer wg.Done()
		UploadToS3(fc.S3Client, file, fc.BucketName, resultChan)
	}()

	wg.Wait()
	result := <-resultChan

	if result == "" || len(result) < 5 {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Upload failed"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"uri": result})
}
