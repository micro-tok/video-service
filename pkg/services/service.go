package services

import (
	"context"
	"os"

	"github.com/gofrs/uuid/v5"
	"github.com/micro-tok/video-service/pkg/cassandra"
	"github.com/micro-tok/video-service/pkg/pb"
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
	cass *cassandra.CassandraService
	s3   *s3.AWSService
}

func NewVideoService(cass *cassandra.CassandraService, s3 *s3.AWSService) Service {
	return &service{
		VideoService: &videoService{
			cass: cass,
			s3:   s3,
		},
	}
}

func (s videoService) UploadVideo(ctx context.Context, req *pb.UploadVideoRequest) (*pb.UploadVideoResponse, error) {
	// Generate UUID
	id, err := uuid.NewV4()
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
	path, err := s.s3.UploadFile(part, id.String()+".mp4")
	if err != nil {
		return nil, err
	}

	// Save to Cassandra
	_, err = s.cass.SaveMetadata(id, req.Title, req.Description, path, req.Tags)
	if err != nil {
		return nil, err
	}

	return &pb.UploadVideoResponse{
		Id: id.String(),
	}, nil
}

func (s videoService) GetVideoMetadata(ctx context.Context, req *pb.GetVideoMetadataRequest) (*pb.GetVideoMetadataResponse, error) {
	id, err := uuid.FromString(req.Id)
	if err != nil {
		return nil, err
	}

	ownerID, title, description, url, tags, err := s.cass.LoadMetadata(id)
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
