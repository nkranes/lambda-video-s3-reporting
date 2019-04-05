package main

import (
	"context"
	"errors"
	"log"
	"time"

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

	resc, errc := make(chan string), make(chan error)

	for i := 0; i < len(folders); i += 1 {
		folder := folders[i]
		go func() {
			folder.SizeInBytes, err = aws.GetFolderSize(sess, config.S3BucketName, folder.Name)

			if err != nil {
				errc <- errors.New(" ERROR: GetFolderSize: " + err.Error())
			} else {
				err = WriteValuesToSql(config.SqlServerConnString, folder)
				if err != nil {
					errc <- errors.New(" ERROR: WriteValuesToSql: " + err.Error())
				}

				resc <- " SUCCESS: Successfully Processed " + folder.Name
			}
		}()
	}

	for i := 0; i < len(folders); i++ {
		select {
		case res := <-resc:
			log.Println(time.Now().Format("2006-01-02 15:04:05") + res)
		case err := <-errc:
			log.Println(time.Now().Format("2006-01-02 15:04:05") + err.Error())
			return err
		}
	}

	return nil
}

func StructifyFolder(objectList *s3.ListObjectsV2Output) []model.Folder {
	var folders []model.Folder
	for _, key := range objectList.CommonPrefixes {
		folder := model.Folder{Name: *key.Prefix}
		folders = append(folders, folder)
	}
	return folders
}

func WriteValuesToSql(sqlServerConnString string, folder model.Folder) (err error) {
	err = services.UpdateTableWithFolderData(sqlServerConnString, folder)
	if err != nil {
		return err
	}

	return nil
}

func main() {
	lambda.Start(MainHandler)
}
