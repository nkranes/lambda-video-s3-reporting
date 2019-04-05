package aws

import (
	"fmt"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

func TestGetBucketTopLevelFoldersOnly(t *testing.T) {
	sess, err := MakeSession()
	if err != nil {
		fmt.Errorf("makeSession error. " + err.Error())
		t.Fail()
	}

	topLevelFolders, err := GetBucketTopLevelFoldersOnly(sess, "frankly-video-content-prod")
	if err != nil {
		fmt.Errorf("GetBucketTopLevelFoldersOnly error. " + err.Error())
		t.Fail()
	}

	for _, key := range topLevelFolders.CommonPrefixes {
		fmt.Println(*key.Prefix)
	}

}

func TestGetFolderSize(t *testing.T) {
	sess, err := MakeSession()
	if err != nil {
		fmt.Errorf("makeSession error. " + err.Error())
		t.Fail()
	}

	folderSize, err := GetFolderSize(sess, "frankly-video-content-prod", "kfol/")
	if err != nil {
		fmt.Errorf("GetFolderSize error. " + err.Error())
		t.Fail()
	}

	if folderSize < 0 {
		t.Fail()
	}
}

func TestObjectSize(t *testing.T) {
	sess, err := MakeSession()
	if err != nil {
		fmt.Errorf("makeSession error. " + err.Error())
		t.Fail()
	}

	svc := s3.New(sess)
	query := &s3.GetObjectInput{
		Bucket: aws.String("frankly-video-content-prod"),
		Key:    aws.String("ap/"),
	}
	resp, err := svc.GetObject(query)

	fmt.Println(resp.ContentLength)
}
