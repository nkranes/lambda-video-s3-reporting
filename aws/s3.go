package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

func GetFolderSize(sess *session.Session, s3BucketName string, s3Subdirectory string) (size int64, err error) {
	query := &s3.ListObjectsV2Input{
		Bucket: &s3BucketName,
		Prefix: aws.String(s3Subdirectory),
	}
	svc := s3.New(sess)

	truncatedListing := true

	for truncatedListing {
		resp, err := svc.ListObjectsV2(query)
		if err != nil {
			fmt.Println("ListObjectsV2 error. " + err.Error())
			return -1, err
		}

		for _, key := range resp.Contents {
			size += *key.Size
		}

		query.ContinuationToken = resp.NextContinuationToken
		truncatedListing = *resp.IsTruncated
	}

	return
}

func GetBucketTopLevelFoldersOnly(sess *session.Session, s3BucketName string) (resp *s3.ListObjectsV2Output, err error) {
	query := &s3.ListObjectsV2Input{
		Bucket:    &s3BucketName,
		Prefix:    aws.String(""),
		Delimiter: aws.String("/"),
	}
	svc := s3.New(sess)

	truncatedListing := true

	for truncatedListing {
		resp, err = svc.ListObjectsV2(query)
		if awsErr, err1 := err.(awserr.Error); err1 {
			switch awsErr.Code() {
			case s3.ErrCodeNoSuchBucket:
				fmt.Println("ListObjectsV2 error. "+s3.ErrCodeNoSuchBucket+": '"+s3BucketName+"'", awsErr.Error())
			default:
				fmt.Println("ListObjectsV2 error. " + awsErr.Error())
			}

			return
		}

		query.ContinuationToken = resp.NextContinuationToken
		truncatedListing = *resp.IsTruncated
	}

	return
}
