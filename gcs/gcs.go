package gcs

import (
	"context"
	"dgn-go-lib/config"
	"dgn-go-lib/consts"
	"dgn-go-lib/general"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"cloud.google.com/go/storage"
	"github.com/labstack/echo/v4"
	"google.golang.org/api/option"
	"google.golang.org/appengine"
)

type GCS struct {
	storageClient *storage.Client
	config        config.BucketConfig
}

type GCSClient interface {
	UploadImage(c echo.Context, imageString, folderName, imageEntityName string) (string, error)
	GetImages(objectName string) (string, error)
}

func DaganganGCSClient(storageClient *storage.Client) *GCS {
	return &GCS{
		storageClient: storageClient,
	}
}

// UploadImageToGCS using base64 image string
func (gcs *GCS) UploadImage(c echo.Context, imageString, folderName, imageEntityName string) (string, error) {
	var bucketName string
	// check is private or public
	bucketName = isPublicOrPrivateBucket(folderName, gcs.config)

	// split image meta data from base64
	base64ImageString := general.Explode(",", imageString)
	reader := base64.NewDecoder(base64.StdEncoding, strings.NewReader(base64ImageString[1]))

	ctx := appengine.NewContext(c.Request())
	storageClient, err := storage.NewClient(ctx, option.WithCredentialsFile(gcs.config.GCPCredentials))
	if err != nil {
		return "A Upload image to GCS failed", err
	}

	// setup image filename
	imageExtension := general.GetMimeFromImage(imageString)
	unixTime := general.GetUnixTimestamp()
	bucketHandle := storageClient.Bucket(bucketName)
	sw := bucketHandle.Object(folderName + "/" + imageEntityName + "-" + strconv.Itoa(int(unixTime)) + imageExtension).NewWriter(ctx)

	if _, err := io.Copy(sw, reader); err != nil {
		return "Upload image to GCS failed", err
	}
	if err := sw.Close(); err != nil {
		return "Upload image to GCS failed", err
	}

	// get image url
	u, err := url.Parse(sw.Attrs().Name)
	if err != nil {
		return "Upload image to GCS failed", err
	}
	imagePath := u.EscapedPath()

	return imagePath, nil
}

func isPublicOrPrivateBucket(folderName string, bucketConfig config.BucketConfig) string {
	parent := filepath.Dir(folderName)
	privateBucketList := general.Explode(",", bucketConfig.PrivateBucketList)
	publicBucketList := general.Explode(",", bucketConfig.PublicBucketList)

	if isExists(parent, privateBucketList) || isExists(folderName, privateBucketList) {
		return bucketConfig.PrivateBucket
	}
	if isExists(parent, publicBucketList) || isExists(folderName, publicBucketList) {
		return bucketConfig.PublicBucket
	}

	return bucketConfig.PublicBucket
}

func (gcs *GCS) GetImages(objectName string) (string, error) {
	var URL = ""
	if strings.Index(objectName, "http") == 1 {
		return objectName, nil
	}
	privateBucketList := general.Explode(",", gcs.config.PrivateBucketList)
	publicBucketList := general.Explode(",", gcs.config.PublicBucketList)

	dirName := filepath.Dir(objectName)
	parent := filepath.Dir(dirName)

	if isExists(dirName, privateBucketList) || isExists(parent, privateBucketList) {
		URL, err := GetPrivateObjectURL(objectName, gcs.config)
		if err != nil {
			log.Println("storage.NewClient: %v", err)
			return "", err
		}
		return URL, nil
	}

	if isExists(dirName, publicBucketList) || isExists(parent, publicBucketList) {
		URL = GetPublicObjectURL(objectName, gcs.config)
		return URL, nil
	}

	return URL, nil
}

// getPrivateObjectSignedURL generates object signed URL with GET method.
func GetPrivateObjectURL(objectName string, bucketConfig config.BucketConfig) (string, error) {
	ctx := context.Background()
	client, err := storage.NewClient(ctx, option.WithCredentialsFile(bucketConfig.GCPCredentials))
	if err != nil {
		log.Println("storage.NewClient: %v", err)
		return "", fmt.Errorf("storage.NewClient: %v", err)
	}
	defer func(client *storage.Client) {
		errs := client.Close()
		if errs != nil {
			log.Println("storage.DeferClient: %v", err)
			return
		}
	}(client)

	opts := &storage.SignedURLOptions{
		Scheme:  storage.SigningSchemeV4,
		Method:  consts.GetRequestMethod,
		Expires: time.Now().Add(1 * time.Minute),
	}

	bucket := client.Bucket(bucketConfig.PrivateBucket)
	url, err := bucket.SignedURL(objectName, opts)
	if err != nil {
		log.Println("Bucket(%q).SignedURL: %v", bucketConfig.PrivateBucket, err)
		return "", fmt.Errorf("Bucket(%q).SignedURL: %v", bucketConfig.PrivateBucket, err)
	}
	return url, nil
}

func GetPublicObjectURL(objectName string, bucketConfig config.BucketConfig) string {
	if strings.Index(objectName, "http") == 1 {
		return objectName
	}
	return bucketConfig.Provider + bucketConfig.PublicBucket + "/" + objectName
}

func isExists(needle string, hayStacks []string) bool {
	for _, v := range hayStacks {
		if strings.Contains(v, needle) {
			return true
		}
	}
	return false
}
