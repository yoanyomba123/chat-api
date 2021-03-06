package storage

import (
	"context"
	"fmt"
	"io/ioutil"

	"github.com/betchi/tracer"
	logger "github.com/betchi/zapper"

	"golang.org/x/oauth2/google"

	storage "google.golang.org/api/storage/v1"
)

type gcsProvider struct {
	ctx                context.Context
	projectId          string
	scope              string
	jwtPath            string
	uploadBucket       string
	uploadDirectory    string
	thumbnailBucket    string
	thumbnailDirectory string
}

var gcsService *storage.Service

func (gp *gcsProvider) Init() error {
	span := tracer.StartSpan(gp.ctx, "Init", "storage")
	defer tracer.Finish(span)

	if gcsService == nil {
		data, err := ioutil.ReadFile(gp.jwtPath)
		if err != nil {
			logger.Error(err.Error())
			tracer.SetError(span, err)
			return err
		}

		conf, err := google.JWTConfigFromJSON(data, gp.scope)
		if err != nil {
			logger.Error(err.Error())
			tracer.SetError(span, err)
			return err
		}

		ctx := context.Background()
		client := conf.Client(ctx)

		service, err := storage.New(client)
		if err != nil {
			logger.Error(err.Error())
			tracer.SetError(span, err)
			return err
		}
		gcsService = service
	}
	return nil
}

func (gp *gcsProvider) Post(assetInfo *AssetInfo) (string, error) {
	span := tracer.StartSpan(gp.ctx, "Post", "storage")
	defer tracer.Finish(span)

	filePath := fmt.Sprintf("%s/%s", gp.uploadDirectory, assetInfo.Filename)
	object := &storage.Object{
		Name: filePath,
	}

	res, err := gcsService.Objects.Insert(gp.uploadBucket, object).Media(assetInfo.Data).Do()
	if err != nil {
		logger.Error(err.Error())
		tracer.SetError(span, err)
		return "", err
	}
	logger.Debug(fmt.Sprintf("name:%s\tselfLink:%s", res.Name, res.SelfLink))

	res, err = gcsService.Objects.Get(gp.uploadBucket, filePath).Do()
	if err != nil {
		logger.Error(err.Error())
		tracer.SetError(span, err)
		return "", err
	}
	logger.Debug(fmt.Sprintf("bucketName:%s\name:%s\tmediaLink:%s", gp.uploadBucket, res.Name, res.MediaLink))

	return res.MediaLink, nil
}

func (gp *gcsProvider) Get(assetInfo *AssetInfo) ([]byte, error) {
	span := tracer.StartSpan(gp.ctx, "Get", "storage")
	defer tracer.Finish(span)

	return nil, nil
}
