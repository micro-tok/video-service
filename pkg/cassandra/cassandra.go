package cassandra

import (
	"log"

	"github.com/gocql/gocql"
	"github.com/gofrs/uuid/v5"
	"github.com/micro-tok/video-service/pkg/config"
)

type CassandraService struct {
	ClusterIP string
	Keyspace  string
}

func NewCassandraService(cfg *config.Config) *CassandraService {
	return &CassandraService{
		ClusterIP: cfg.CassandraClusterIP,
		Keyspace:  cfg.CassandraKeyspace,
	}
}

func (s CassandraService) SaveMetadata(ownerID uuid.UUID, title string, description string, url string, tags []string) (uuid.UUID, error) {
	cluster := gocql.NewCluster(s.ClusterIP)
	cluster.Keyspace = s.Keyspace
	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatalf("Failed to connect to Cassandra: %v", err)
		return uuid.UUID{}, err
	}
	defer session.Close()

	videoID, err := uuid.NewV4()
	if err != nil {
		return uuid.UUID{}, err
	}

	if err := session.Query(`INSERT INTO videos (id, ownerId, title, description, url, tags) VALUES (?, ?, ?, ?, ?, ?)`, videoID, ownerID, title, description, url, tags).Exec(); err != nil {
		log.Fatalf("Failed to insert into Cassandra: %v", err)
		return uuid.UUID{}, err
	}

	return videoID, nil
}

func (s CassandraService) LoadMetadata(videoID uuid.UUID) (uuid.UUID, string, string, string, []string, error) {
	cluster := gocql.NewCluster(s.ClusterIP)
	cluster.Keyspace = s.Keyspace
	session, err := cluster.CreateSession()
	if err != nil {
		log.Fatalf("Failed to connect to Cassandra: %v", err)
		return uuid.UUID{}, "", "", "", nil, err
	}
	defer session.Close()

	var ownerID uuid.UUID
	var title string
	var description string
	var url string
	var tags []string

	if err := session.Query(`SELECT ownerId, title, description, url, tags FROM videos WHERE id = ?`, videoID).Scan(&ownerID, &title, &description, &url, &tags); err != nil {
		log.Fatalf("Failed to select from Cassandra: %v", err)
		return uuid.UUID{}, "", "", "", nil, err
	}

	return ownerID, title, description, url, tags, nil
}
