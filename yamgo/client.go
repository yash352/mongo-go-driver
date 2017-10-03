package yamgo

import (
	"context"
	"time"

	"github.com/10gen/mongo-go-driver/yamgo/connstring"
	"github.com/10gen/mongo-go-driver/yamgo/private/cluster"
	"github.com/10gen/mongo-go-driver/yamgo/private/ops"
	"github.com/10gen/mongo-go-driver/yamgo/readpref"
)

const defaultLocalThreshold = 15 * time.Millisecond

// Client performs operations on a given cluster.
type Client struct {
	cluster        *cluster.Cluster
	connString     connstring.ConnString
	localThreshold time.Duration
	readPreference *readpref.ReadPref
}

// NewClient creates a new client to connect to a cluster specified by the uri.
func NewClient(uri string) (*Client, error) {
	cs, err := connstring.Parse(uri)
	if err != nil {
		return nil, err
	}

	return NewClientFromConnString(cs)
}

// NewClientFromConnString creates a new client to connect to a cluster specified by the connection string.
func NewClientFromConnString(cs connstring.ConnString) (*Client, error) {
	clst, err := cluster.New(cluster.WithConnString(cs))
	if err != nil {
		return nil, err
	}

	// TODO GODRIVER-92: Allow custom localThreshold
	client := &Client{
		cluster:        clst,
		connString:     cs,
		localThreshold: defaultLocalThreshold,
		readPreference: readpref.Primary(),
	}

	return client, nil
}

// Database returns a handle for a given database.
func (client *Client) Database(name string, options ...DatabaseOption) *Database {
	return newDatabase(client, name, options...)
}

// ConnectionString returns the connection string of the cluster the client is connected to.
func (client *Client) ConnectionString() string {
	return client.connString.Original
}

func (client *Client) selectServer(ctx context.Context, selector cluster.ServerSelector,
	readPref *readpref.ReadPref) (*ops.SelectedServer, error) {

	s, err := client.cluster.SelectServer(ctx, selector, readPref)
	if err != nil {
		return nil, err
	}

	return s, nil
}