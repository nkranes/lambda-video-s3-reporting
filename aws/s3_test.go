package aws

import (
	"fmt"
	"testing"
)


func TestGetBucketObjectsInvalid(t *testing.T) {
	sess, err := MakeSession()
	if err != nil {
		fmt.Errorf("makeSession error. " + err.Error())
		t.Fail()
	}

	_, err = GetBucketObjects(sess, "frankly-video-content-prodXXX", "")
	if err != nil {
		fmt.Println("Success")
	} else {
		t.Fail()
	}
}

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

func TestGetBucketObjects(t *testing.T) {
	sess, err := MakeSession()
	if err != nil {
		fmt.Errorf("makeSession error. " + err.Error())
		t.Fail()
	}

	objects, err := GetBucketObjects(sess, "frankly-video-content-prod", "ap/")
	if err != nil {
		fmt.Errorf("GetBucketObjects error. " + err.Error())
		t.Fail()
	}

	fmt.Println(objects.Contents)

	for _, key := range objects.Contents {
		fmt.Println(*key.Size)
	}
}

