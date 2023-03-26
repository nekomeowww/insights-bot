package pinecone

import (
	"context"
	"crypto/tls"
	"fmt"
	"time"

	"github.com/pinecone-io/go-pinecone/pinecone_grpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

type Client struct {
	pinecone_grpc.VectorServiceClient
	conn *grpc.ClientConn
}

func NewClient(
	index string,
	project string,
	environment string,
	apiKey string,
) (*Client, error) {
	if index == "" {
		return nil, fmt.Errorf("Pinecone index name is required")
	}
	if project == "" {
		return nil, fmt.Errorf("Pinecone project name is required")
	}
	if environment == "" {
		return nil, fmt.Errorf("Pinecone environment name is required")
	}

	tlsConfig := &tls.Config{}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	ctx = metadata.AppendToOutgoingContext(ctx, "api-key", apiKey)
	pineconeTarget := fmt.Sprintf("%s-%s.svc.%s.pinecone.io:443", index, project, environment)
	conn, err := grpc.DialContext(
		ctx,
		pineconeTarget,
		grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)),
		grpc.WithAuthority(pineconeTarget),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, err
	}

	return &Client{
		VectorServiceClient: pinecone_grpc.NewVectorServiceClient(conn),
		conn:                conn,
	}, nil
}

func (c *Client) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}
