package storage

import (
	"fmt"

	"bytes"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
	"github.com/go-kit/kit/log"
	"github.com/monax/hoard/logging"
	"github.com/monax/hoard/logging/structure"
)

type s3Store struct {
	awsS3           *s3.S3
	awsUploader     *s3manager.Uploader
	awsDownloader   *s3manager.Downloader
	s3Bucket        string
	s3Prefix        string
	addressEncoding AddressEncoding
	logger          log.Logger
}

const NotFoundCode = "NotFound"

func NewS3Store(s3Bucket, s3Prefix string, addressEncoding AddressEncoding,
	awsConfig *aws.Config, logger log.Logger) (*s3Store, error) {

	if awsConfig == nil {
		awsConfig = aws.NewConfig()
	}

	if logger == nil {
		logger = log.NewNopLogger()
	}

	awsSession, err := Session(awsConfig)
	if err != nil {
		return nil, err
	}
	s3s := &s3Store{
		awsS3:           s3.New(awsSession),
		awsUploader:     s3manager.NewUploader(awsSession),
		awsDownloader:   s3manager.NewDownloader(awsSession),
		s3Bucket:        s3Bucket,
		s3Prefix:        s3Prefix,
		addressEncoding: addressEncoding,
		logger: logging.TraceLogger(log.With(logger,
			structure.ComponentKey, "storage")),
	}
	s3s.logger = log.With(s3s.logger, "store_name", s3s.Name())
	return s3s, nil
}

func Session(awsConfig *aws.Config) (*session.Session, error) {
	return session.NewSessionWithOptions(session.Options{
		Config:            *awsConfig,
		SharedConfigState: session.SharedConfigEnable,
	})
}

func (s3s *s3Store) Put(address []byte, data []byte) ([]byte, error) {
	// Should be threadsafe
	output, err := s3s.awsUploader.Upload(&s3manager.UploadInput{
		Bucket: &s3s.s3Bucket,
		Key:    aws.String(s3s.Key(address)),
		Body:   bytes.NewReader(data),
	})
	if err != nil {
		return address, err
	}
	s3s.logger.Log("method", "Put",
		"location", output.Location,
		"encoded_address", s3s.encode(address),
		"version_id", output.VersionID,
		"upload_id", output.UploadID)
	return address, err
}

func (s3s *s3Store) Get(address []byte) ([]byte, error) {
	buf := &aws.WriteAtBuffer{}
	n, err := s3s.awsDownloader.Download(buf, &s3.GetObjectInput{
		Bucket: &s3s.s3Bucket,
		Key:    aws.String(s3s.Key(address)),
	})
	s3s.logger.Log("method", "Get",
		"encoded_address", s3s.encode(address),
		"downloaded_bytes", n)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (s3s *s3Store) Stat(address []byte) (*StatInfo, error) {
	output, err := s3s.awsS3.HeadObject(&s3.HeadObjectInput{
		Bucket: &s3s.s3Bucket,
		Key:    aws.String(s3s.Key(address)),
	})
	if err != nil {
		s3err, ok := err.(awserr.Error)
		if ok && s3err.Code() == NotFoundCode {
			return &StatInfo{
				Exists: false,
			}, nil
		}
		return nil, err
	}
	s3s.logger.Log("method", "Stat",
		"encoded_address", s3s.encode(address),
		"version_id", output.VersionId,
		"etag", output.ETag)
	return &StatInfo{
		Exists: true,
		Size:   uint64(*output.ContentLength),
	}, nil
}

func (s3s *s3Store) Location(address []byte) string {
	return fmt.Sprintf("https://%s.s3.amazonaws.com/%s", s3s.s3Bucket,
		s3s.Key(address))
}

func (s3s *s3Store) Key(address []byte) string {
	return fmt.Sprintf("%s/%s", s3s.s3Prefix, s3s.encode(address))
}

func (s3s *s3Store) Name() string {
	return fmt.Sprintf("s3Store[bucket=%s,prefix=%s]", s3s.s3Bucket,
		s3s.s3Prefix)
}

func (s3s *s3Store) encode(address []byte) string {
	return s3s.addressEncoding.EncodeToString(address)
}
