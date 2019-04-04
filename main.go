package main

import (
	"context"
	"errors"
	"log"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/kelseyhightower/envconfig"

	"github.com/franklyinc/frankly-lambda-video-s3-reporting/aws"
	"github.com/franklyinc/frankly-lambda-video-s3-reporting/model"
	"github.com/franklyinc/frankly-lambda-video-s3-reporting/services"
)

func MainHandler(ctx context.Context) error {
	var config model.Config
	err := envconfig.Process("", &config)
	if err != nil {
		return errors.New("ERROR: failed to process the environment variables properly: " + err.Error())
	}
	log.Println(config)

	if config.Pause {
		return nil
	}

	sess, err := aws.MakeSession()
	if err != nil {
		return errors.New("ERROR: makeSession: " + err.Error())
	}

	topLevelFolders, err := aws.GetBucketTopLevelFoldersOnly(sess, config.S3BucketName)
	if err != nil {
		return errors.New("ERROR: GetBucketTopLevelFoldersOnly: " + err.Error())
	}

	folders := StructifyFolder(topLevelFolders)

	for i := 0; i < len(folders); i += 1 {
		folders[i].SizeInBytes, err = aws.GetFolderSize(sess, config.S3BucketName, folders[i].Name)
		if err != nil {
			return errors.New("ERROR: GetFolderSize: " + err.Error())
		}
	}

	err = WriteValuesToSql(config.SqlServerConnString, folders)
	if err != nil {
		return errors.New("ERROR: WriteValuesToSql: " + err.Error())
	}

	return nil
}

func StructifyFolder(objectList *s3.ListObjectsV2Output) ([]model.Folder){
	var folders []model.Folder
	for _, key := range objectList.CommonPrefixes {
		folder := model.Folder{Name: *key.Prefix }
		folders = append(folders, folder)
	}
	return folders
}

func WriteValuesToSql(sqlServerConnString string, folders []model.Folder) (err error) {
	for _, folder := range folders {
		err = services.UpdateTableWithFolderData(sqlServerConnString, folder)
		if err != nil {
			return err
		}
	}

	return nil
}

func main() {
	lambda.Start(MainHandler)
}
