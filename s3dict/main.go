package s3dict

// This file contains the code to query the dictionary file in AWS S3
// using the AWS SDK for Go.

import (
	"bytes"
	"context"
	"encoding/binary"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/harshjoeyit/word-dict/dict"
)

type S3Dict struct {
	// key is dictonary file's relative path in S3 bucket
	key string
	s3b *S3Bucket
	// in-memory index of the dictionary
	index map[string]dict.IndexEntry
}

func New() (*S3Dict, error) {
	s3b, err := NewS3Bucket()
	if err != nil {
		return nil, fmt.Errorf("failed to create S3 bucket client: %v", err)
	}

	key := os.Getenv("DICT_KEY")

	// Read the index from key file
	index, err := readIndex(s3b, key)
	if err != nil {
		log.Fatalf("failed to read index: %v", err)
	}

	s3d := &S3Dict{
		key:   key,
		s3b:   s3b,
		index: index,
	}

	return s3d, nil
}

// QueryWord queries the dictionary for a word and returns its definition
func (d *S3Dict) QueryWord(word string) (string, bool) {
	// Find the word in the index
	idxe, ok := d.index[word]
	if !ok {
		log.Printf("Error: '%s' word not found", word)
		return "", false
	}

	// Read the definition - <word>,<definition>
	defOffsetSt := idxe.Offset + int64(len(idxe.Word)+1) // +1 for the comma
	defOffsetEn := defOffsetSt + int64(idxe.DefSize) - 1

	data, err := d.s3b.GetObjectByteRange(d.key, defOffsetSt, defOffsetEn) // +1 for the comma
	if err != nil {
		log.Printf("error reading definition: %v", err)
		return "", false
	}

	return string(data), true
}

type S3Bucket struct {
	bucketName string
	client     *s3.Client
}

func NewS3Bucket() (*S3Bucket, error) {
	// Get bucket credentials from environment variables
	// These should be set in the .env file
	bucketName := os.Getenv("S3_BUCKET_NAME")
	region := os.Getenv("AWS_REGION")
	accessKeyID := os.Getenv("AWS_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("AWS_SECRET_ACCESS_KEY")

	// Config
	cfg, err := config.LoadDefaultConfig(
		context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			accessKeyID, secretAccessKey, "",
		)),
	)
	if err != nil {
		return nil, fmt.Errorf("unable to load SDK config, %v", err)
	}

	s3b := &S3Bucket{
		bucketName: bucketName,
		client:     s3.NewFromConfig(cfg),
	}

	return s3b, nil
}

func readIndex(s3b *S3Bucket, key string) (map[string]dict.IndexEntry, error) {
	// Read the constant sized index header - 8 bytes
	const indexHeaderSize int64 = 8
	data, err := s3b.GetObjectByteRange(key, 0, indexHeaderSize-1)
	if err != nil {
		return nil, fmt.Errorf("unable to get index header, %v", err)
	}

	// Convert byte slice to int64
	indexSize := int64(binary.BigEndian.Uint64(data))

	log.Println("Index size:", indexSize)

	// Read the index
	indexOffsetSt := indexHeaderSize
	indexOffsetEn := indexOffsetSt + indexSize - 1

	data, err = s3b.GetObjectByteRange(key, indexOffsetSt, indexOffsetEn)
	if err != nil {
		return nil, fmt.Errorf("unable to get index data, %v", err)
	}

	index := make(map[string]dict.IndexEntry)

	// Convert byte slice to string and split by new line
	// Each line is of the format <word>:<offset>:<defSize>
	lines := strings.Split(string(data), "\n")

	for _, line := range lines {
		if line == "" {
			continue
		}

		parts := strings.Split(line, ":")
		if len(parts) != 3 {
			continue
		}

		word := parts[0]
		offset := int64(binary.BigEndian.Uint64([]byte(parts[1])))
		defSize := int16(binary.BigEndian.Uint16([]byte(parts[2])))

		index[word] = dict.IndexEntry{
			Word:    word,
			Offset:  offset,
			DefSize: defSize,
		}
	}

	log.Println("Total index entries:", len(index))

	return index, nil
}

// GetObjectByteRange retrieves an object from S3 for byte range [rangeSt, rangeEn]
// (both inclusive) and returns the data as a byte slice.
func (s3b *S3Bucket) GetObjectByteRange(key string, rangeSt, rangeEn int64) ([]byte, error) {
	// Create the GetObjectInput with Range parameter
	input := &s3.GetObjectInput{
		Bucket: aws.String(s3b.bucketName),
		Key:    aws.String(key),
		Range:  aws.String(fmt.Sprintf("bytes=%d-%d", rangeSt, rangeEn)),
	}

	// Perform the request
	result, err := s3b.client.GetObject(context.TODO(), input)
	if err != nil {
		return nil, fmt.Errorf("failed to get object from s3, %v", err)
	}
	defer result.Body.Close()

	// Get the bytes from the response body
	var buf bytes.Buffer
	if _, err := buf.ReadFrom(result.Body); err != nil {
		return nil, fmt.Errorf("unable to read object data, %v", err)
	}

	data := buf.Bytes()
	if len(data) == 0 {
		return nil, fmt.Errorf("no data returned for range %d-%d", rangeSt, rangeEn)
	}

	fmt.Printf("Downloaded bytes %d-%d from %s/%s\n", rangeSt, rangeEn, s3b.bucketName, key)

	return data, nil
}
