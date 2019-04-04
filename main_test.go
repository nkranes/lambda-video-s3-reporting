package main

import (
	"fmt"
	"testing"

	"github.com/franklyinc/frankly-lambda-video-s3-reporting/model"
)

func TestMainHandler(t *testing.T) {
	MainHandler(nil)
}

func TestWriteValuesToSql(t *testing.T) {
	var folders = []model.Folder {
		{
			Name: "AP",
			SizeInBytes:123456,
		},
	}

	err := WriteValuesToSql("XX", folders)
	if err != nil {
		t.Fail()
	}
}

func TestFloatingPointMath(t *testing.T) {
	var sizeInMegaBytes float64
	sizeInMegaBytes = 7031781042 / 1024.0

	var temp1 float64
	temp1 = 1.23

	fmt.Println(sizeInMegaBytes)
	fmt.Println(temp1)
}
