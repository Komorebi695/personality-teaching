package Cos

import (
	"context"
	"github.com/tencentyun/cos-go-sdk-v5"
	"net/http"
	"net/url"
)

func NewCosClient(secretId, secretKey, region, bucket string)(*cos.Client, error){

	cosURL,err := url.Parse("https://" + bucket + ".cos." + region + ".myqcloud.com")
	if err !=nil{
		return nil,err
	}

	client := cos.NewClient(&cos.BaseURL{
		BucketURL: cosURL,
	}, &http.Client{
		Transport: &cos.AuthorizationTransport{
			SecretID:  secretId,
			SecretKey: secretKey,
		},
	})
	_, err = client.Bucket.Head(context.Background(), &cos.BucketHeadOptions{})
	if err != nil {
		return nil, err
	}
	return client, nil
}
