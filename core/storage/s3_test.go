// +build integration

package storage

// func TestS3Store(t *testing.T) {
// 	bucket := "monax-hoard-test"
// 	prefix := "TestS3Store/"
// 	deletePrefix(bucket, prefix)
// 	s3s, err := NewS3Store(bucket, prefix, base32.StdEncoding, nil, nil)
// 	assert.NoError(t, err)
// 	testStore(t, s3s)
// }

// func deletePrefix(bucket, prefix string) {
// 	deleter := s3manager.NewBatchDelete(session.Must(Session(aws.NewConfig())))
// 	err := deleter.Delete(context.Background(),
// 		s3manager.NewDeleteListIterator(deleter.Client,
// 			&s3.ListObjectsInput{Bucket: &bucket, Prefix: &prefix}))
// 	if err != nil {
// 		panic(err)
// 	}
// }
