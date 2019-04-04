package model

type Config struct {
	SqlServerConnString string  `envconfig:"SQL_SERVER_CONN_STRING"`
	S3BucketName        string  `envconfig:"S3_BUCKET_NAME"`
	Pause               bool    `envconfig:"PAUSE"`
}
