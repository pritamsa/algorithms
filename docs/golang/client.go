package s3

import (
	"io/ioutil"
	"log"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	s32 "github.com/aws/aws-sdk-go/service/s3"
)

type Client interface {
	Get() ([]byte, error)
}

var s, _ = session.NewSession(&aws.Config{Region: aws.String("us-west-2")})

type SecureS3Client struct {
	client *s32.S3
	uri    string
}

func NewSecureS3Client(uri string) Client {
	s3 := s32.New(s)
	return &SecureS3Client{s3, uri}
}

func (s *SecureS3Client) Get() (object []byte, err error) {
	// in k8s
	if !strings.HasPrefix(s.uri, "keys/") {
		return []byte(s.uri), nil
	}

	for i := 1; i <= 10; i++ {
		log.Println("s3.Get()", s.uri, "attempt:", i)
		output, err := s.client.GetObject(
			&s32.GetObjectInput{
				Bucket: aws.String("p3-nonprod"),
				Key:    aws.String(s.uri),
			})

		if err != nil {
			log.Println(err)
			time.Sleep(1 * time.Second)
			continue
		}

		defer output.Body.Close()
		return ioutil.ReadAll(output.Body)
	}
	return nil, err
}
