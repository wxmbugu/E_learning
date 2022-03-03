package aws

import (
	"log"

	"github.com/E_learning/util"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

func ConnectAws() *session.Session {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal(err)
	}
	Accesskeyid := config.Awsaccesskey
	Secretkeyaccess := config.Awssecretkey
	Region := config.Awsregion

	session, err := session.NewSession(&aws.Config{
		Region: &Region,
		Credentials: credentials.NewStaticCredentials(
			Accesskeyid,
			Secretkeyaccess,
			"",
		),
	})
	if err != nil {
		log.Fatal(err)
	}
	return session
}
