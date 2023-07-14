package cloudinaryutil

import (
	"bytes"
	"context"
	"github.com/cloudinary/cloudinary-go/v2"
	"github.com/cloudinary/cloudinary-go/v2/api/uploader"
	"io/ioutil"
	"mime/multipart"
)

type CloudinaryClient struct {
	client *cloudinary.Cloudinary
}

func NewCloudinaryClient(cloudName, apiKey, apiSecret string) *CloudinaryClient {
	client, _ := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
	return &CloudinaryClient{
		client: client,
	}
}

func (c *CloudinaryClient) UploadImage(ctx context.Context, imageFile *multipart.FileHeader) (string, error) {
	src, err := imageFile.Open()
	if err != nil {
		return "", err
	}
	defer src.Close()

	imageData, err := ioutil.ReadAll(src)
	if err != nil {
		return "", err
	}

	uploadResult, err := c.client.Upload.Upload(ctx, bytes.NewReader(imageData), uploader.UploadParams{})
	if err != nil {
		return "", err
	}

	return uploadResult.SecureURL, nil
}
