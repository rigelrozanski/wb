package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

func restore() {

	bucket, sess := loadOpenSess()

	// Get the list of objects
	svc := s3.New(sess)
	resp, err := svc.ListObjects(&s3.ListObjectsInput{Bucket: aws.String(bucket)})
	if err != nil {
		exitErrorf("Unable to list items in bucket %q, %v", bucket, err)
	}

	//loop through the items and download
	for _, o := range resp.Contents {

		name := *o.Key
		boardPath, err := getWbPath(name)
		if err != nil {
			exitErrorf("Unable to get path %q, %v", err)
		}

		file, err := os.Create(boardPath)

		if err != nil {
			exitErrorf("Unable to open file %q, %v", err)
		}

		defer file.Close()

		downloader := s3manager.NewDownloader(sess)

		numBytes, err := downloader.Download(file,
			&s3.GetObjectInput{
				Bucket: aws.String(bucket),
				Key:    aws.String(name),
			})

		if err != nil {
			exitErrorf("Unable to download item %q, %v", name, err)
		}

		fmt.Println("Downloaded", file.Name(), numBytes, "bytes")
	}
}

func loadOpenSess() (bucket string, sess *session.Session) {
	//unmarshal the settings
	keypath, err := getKeyPath()
	if err != nil {
		panic(err)
	}
	settingsFile, err := ioutil.ReadFile(keypath)
	if err != nil {
		panic(err)
	}
	var a map[string]interface{}
	json.Unmarshal(settingsFile, &a)
	bucket = a["bucket"].(string)
	id := a["aws_id"].(string)
	secret := a["aws_secret"].(string)
	region := a["aws_region"].(string)

	//make the AWS session
	sess = session.Must(session.NewSession(&aws.Config{
		Region:      aws.String(region),
		Credentials: credentials.NewStaticCredentials(id, secret, ""),
	}))

	return
}

func backup() {
	bucket, sess := loadOpenSess()

	////////////////////////////////////////////////////////////////
	// Empty the bucket
	////////////////////////////////////////////////////////////////

	// Get the list of objects
	svc := s3.New(sess)
	resp, err := svc.ListObjects(&s3.ListObjectsInput{Bucket: aws.String(bucket)})
	if err != nil {
		exitErrorf("Unable to list items in bucket %q, %v", bucket, err)
	}
	numObjs := len(resp.Contents)

	// Create Delete object with slots for the objects to delete
	var items s3.Delete
	var objs = make([]*s3.ObjectIdentifier, numObjs)

	for i, o := range resp.Contents {
		// Add objects from command line to array
		objs[i] = &s3.ObjectIdentifier{Key: aws.String(*o.Key)}
	}

	// Add list of objects to delete to Delete object
	items.SetObjects(objs)

	// Delete the items
	_, err = svc.DeleteObjects(&s3.DeleteObjectsInput{Bucket: &bucket, Delete: &items})
	if err != nil {
		exitErrorf("Unable to delete objects from bucket %q, %v", bucket, err)
	}
	fmt.Println("Deleted", numObjs, "object(s) from bucket", bucket)

	////////////////////////////////////////////////////////////////
	// Upload all the files
	////////////////////////////////////////////////////////////////

	// Setup the S3 Upload Manager. Also see the SDK doc for the Upload Manager
	// for more information on configuring part size, and concurrency.
	//
	// http://docs.aws.amazon.com/sdk-for-go/api/service/s3/s3manager/#NewUploader
	uploader := s3manager.NewUploader(sess)

	boardPath, err := getWbPath("")
	if err != nil {
		fmt.Println(err)
		return
	}

	visit := func(filepath string, f os.FileInfo, err error) error {
		// Open the file
		file, err := os.Open(filepath)
		if err != nil {
			return fmt.Errorf("Unable to open file %v", err)
		}
		defer file.Close()

		// Create the key name for the bucket
		basePath := path.Base(filepath)
		name := strings.Replace(basePath, boardsDir, "", 1) //remove the boards dir
		if len(name) > 0 {
			// Upload to bucket with the key being the same as the filename.
			_, err = uploader.Upload(&s3manager.UploadInput{
				Bucket: aws.String(bucket),
				Key:    aws.String(name),
				Body:   file,
			})
			if err != nil {
				exitErrorf("Unable to upload %q to %q, %v", filepath, bucket, err)
			}
			fmt.Printf("Successfully uploaded %q to %q\n", filepath, bucket)
		}
		return nil
	}
	filepath.Walk(boardPath, visit)
}

func exitErrorf(msg string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, msg+"\n", args...)
	os.Exit(1)
}
