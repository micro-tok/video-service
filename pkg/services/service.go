package services

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"

	"github.com/gofrs/uuid/v5"
	"github.com/micro-tok/video-service/pkg/cassandra"
	"github.com/micro-tok/video-service/pkg/pb"
	"github.com/micro-tok/video-service/pkg/redis"
	"github.com/micro-tok/video-service/pkg/s3"
)

type Service interface {
	VideoService
}

type service struct {
	VideoService
}

type VideoService interface {
	UploadVideo(ctx context.Context, req *pb.UploadVideoRequest) (*pb.UploadVideoResponse, error)
	GetVideoMetadata(ctx context.Context, req *pb.GetVideoMetadataRequest) (*pb.GetVideoMetadataResponse, error)
}

type videoService struct {
	cass  *cassandra.CassandraService
	s3    *s3.AWSService
	redis *redis.RedisService
}

func NewVideoService(cass *cassandra.CassandraService, s3 *s3.AWSService, redis *redis.RedisService) Service {
	return &service{
		VideoService: &videoService{
			cass:  cass,
			s3:    s3,
			redis: redis,
		},
	}
}

func (s videoService) UploadVideo(ctx context.Context, req *pb.UploadVideoRequest) (*pb.UploadVideoResponse, error) {
	// Generate UUID
	id := uuid.FromStringOrNil(req.OwnerId)
	if id == uuid.Nil {
		fmt.Println("Invalid owner id: ", req.OwnerId)
		return nil, errors.New("invalid owner id")
	}

	videoID, err := uuid.NewV4()
	if err != nil {
		return nil, err
	}

	tmpfile, err := os.CreateTemp("", "video-*.mp4")
	if err != nil {
		return nil, err
	}
	defer os.Remove(tmpfile.Name())

	// Write to temp file
	_, err = tmpfile.Write(req.Video)
	if err != nil {
		return nil, err
	}

	// Open temp file
	part, err := os.Open(tmpfile.Name())
	if err != nil {
		return nil, err
	}

	// Upload to S3
	path, err := s.s3.UploadFile(part, videoID.String()+".mp4")
	if err != nil {
		return nil, err
	}

	// Save to Cassandra
	_, err = s.cass.SaveMetadata(videoID, id, req.Title, req.Description, path, req.Tags)
	if err != nil {
		return nil, err
	}

	return &pb.UploadVideoResponse{
		Id: videoID.String(),
	}, nil
}

func (s videoService) GetVideoMetadata(ctx context.Context, req *pb.GetVideoMetadataRequest) (*pb.GetVideoMetadataResponse, error) {
	val, _ := s.redis.Get(req.Id)

	if val != "" {
		//unmarshal the value
		var res pb.GetVideoMetadataResponse
		err := json.Unmarshal([]byte(val), &res)
		if err != nil {
			return nil, err
		}

		return &res, nil
	}

	id, err := uuid.FromString(req.Id)
	if err != nil {
		return nil, err
	}

	ownerID, title, description, url, tags, err := s.cass.LoadMetadata(id)
	if err != nil {
		return nil, err
	}

	res := &pb.GetVideoMetadataResponse{
		Id:          id.String(),
		OwnerId:     ownerID.String(),
		Title:       title,
		Description: description,
		Url:         url,
		Tags:        tags,
	}

	//marshal the value
	jsonData, err := json.Marshal(res)
	if err != nil {
		return nil, err
	}

	//set the value in redis
	err = s.redis.Set(req.Id, jsonData)
	if err != nil {
		return nil, err
	}

	return &pb.GetVideoMetadataResponse{
		Id:          id.String(),
		OwnerId:     ownerID.String(),
		Title:       title,
		Description: description,
		Url:         url,
		Tags:        tags,
	}, nil
}
