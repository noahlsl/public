package obsx

import (
	"github.com/huaweicloud/huaweicloud-sdk-go-obs/obs"
)

type Conn struct {
	bucketName string
	objectKey  string
	location   string
	obsClient  *obs.ObsClient
}

func (c *Cfg) NewClient() *Conn {

	obsClient, err := obs.New(c.Ak, c.Sk, c.Endpoint)
	if err != nil {
		panic(err)
	}

	conn := newConn(obsClient, c.BucketName, c.ObjectKey, c.Location)
	input := &obs.CreateBucketInput{}
	input.Bucket = c.BucketName
	input.Location = c.Location
	_, err = conn.obsClient.CreateBucket(input)
	if err != nil {
		panic(err)
	}

	return conn
}

func (c *Cfg) MustObsClient() *obs.ObsClient {
	obsClient, err := obs.New(c.Ak, c.Sk, c.Endpoint)
	if err != nil {
		panic(err)
	}
	return obsClient
}

func newConn(obsClient *obs.ObsClient, bucketName, objectKey, location string) *Conn {
	return &Conn{
		bucketName: bucketName,
		objectKey:  objectKey,
		location:   location,
		obsClient:  obsClient,
	}
}
