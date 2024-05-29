package repository

import (
	"context"
	"fmt"
	"log"
	"memorizor/services/account/util"
	"mime/multipart"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gofrs/uuid"
)

type sProfileImageRepositoryAWS struct {
	client        *s3.Client
	bucketName    string
	bucketURLRoot string
}

func NewSProfileImageRepositoryAWS(client *s3.Client, bucketName string, bucketURLRoot string) IProfileImageRepository {
	return &sProfileImageRepositoryAWS{
		client:        client,
		bucketName:    bucketName,
		bucketURLRoot: bucketURLRoot,
	}
}

func (r *sProfileImageRepositoryAWS) Update(userID uuid.UUID, imageFile multipart.File, imageType string) (imageURL string, err error) {
	imageURL = fmt.Sprintf("profile_images/%s.%s", userID.String(), imageType)
	log.Println("Begin to upload an image")
	_, err = r.client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(r.bucketName),
		Key:    aws.String(imageURL),
		Body:   imageFile,
	})
	log.Println("Upload image completed")
	if err != nil {
		return "", util.NewInternal("Unable to update the profile image")
	}
	return r.bucketURLRoot + imageURL, err
}
