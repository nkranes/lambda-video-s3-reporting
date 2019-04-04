package aws

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

func MakeSession() (sess *session.Session, err error)  {
	sess, err = session.NewSessionWithOptions(session.Options{
		Config: aws.Config{
			Region:      aws.String("us-east-1"),
			Credentials: credentials.NewEnvCredentials(),
		},
	})
	if err != nil {
		fmt.Println("failed to create session,", err)
		return
	}

	return
}

