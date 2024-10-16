package server

import (
	"context"
	"database/sql"
	"log"

	pb "github.com/spodolaks/qbility-backend/generated" // Path to generated Protobuf Go files
)

// TopicServer implements the gRPC server defined in the proto file.
type TopicServer struct {
	pb.UnimplementedTopicServiceServer         // Embed to have forward compatibility
	db                                 *sql.DB // Database connection pool
}

// NewTopicServer creates a new TopicServer with a DB connection
func NewTopicServer(db *sql.DB) *TopicServer {
	return &TopicServer{db: db}
}

// GetTopics fetches the list of topics from the MySQL database
func (s *TopicServer) GetTopics(ctx context.Context, req *pb.Empty) (*pb.TopicList, error) {
	rows, err := s.db.Query("SELECT id, title FROM topics")
	if err != nil {
		log.Printf("Error fetching topics: %v", err)
		return nil, err
	}
	defer rows.Close()

	var topics []*pb.Topic
	for rows.Next() {
		var topic pb.Topic
		if err := rows.Scan(&topic.Id, &topic.Title); err != nil {
			return nil, err
		}
		topics = append(topics, &topic)
	}

	return &pb.TopicList{Topics: topics}, nil
}

// GetTopicContent fetches the HTML content of a specific topic
func (s *TopicServer) GetTopicContent(ctx context.Context, req *pb.TopicRequest) (*pb.TopicContent, error) {
	var content string
	err := s.db.QueryRow("SELECT html_content FROM topics WHERE id = ?", req.TopicId).Scan(&content)
	if err != nil {
		log.Printf("Error fetching topic content: %v", err)
		return nil, err
	}

	return &pb.TopicContent{TopicId: req.TopicId, HtmlContent: content}, nil
}
