package routeguide

import (
	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"grpcdemo/pkg/client"
	"grpcdemo/proto/routeguide"
	"io"
	"log"
	"math/rand"
	"time"
)

type Client struct {
	conn   *grpc.ClientConn
	client routeguide.RouteGuideClient
	config *client.RPCClient
}

func WithConsul(consulAddr string) client.RPCClientOption {
	return client.WithTarget("consul://" + consulAddr + "/routeguide-Routeguide")
}

func NewClient(opts ...client.RPCClientOption) (*Client, error) {
	c := Client{config: &client.RPCClient{}}

	for _, opt := range opts {
		opt(c.config)
	}

	conn, err := c.config.Dial()
	if err != nil {
		return nil, err
	}
	return &Client{conn: conn, client: routeguide.NewRouteGuideClient(conn)}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

// PrintFeature gets the feature for the given point.
func (c *Client) PrintFeature(point *routeguide.Point) {
	log.Printf("Getting feature for point (%d, %d)", point.Latitude, point.Longitude)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	feature, err := c.client.GetFeature(ctx, point)
	if err != nil {
		log.Fatalf("%v.GetFeatures(_) = _, %v: ", c.client, err)
	}
	log.Println(feature)
}

// PrintFeatures lists all the features within the given bounding Rectangle.
func (c *Client) PrintFeatures(rect *routeguide.Rectangle) {
	log.Printf("Looking for features within %v", rect)
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	stream, err := c.client.ListFeatures(ctx, rect)
	if err != nil {
		log.Fatalf("%v.ListFeatures(_) = _, %v", c.client, err)
	}
	for {
		feature, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("%v.ListFeatures(_) = _, %v", c.client, err)
		}
		log.Println(feature)
	}
}

// RunRecordRoute sends a sequence of points to server and expects to get a RouteSummary from server.
func (c *Client) RunRecordRoute() {
	// Create a random number of random points
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	pointCount := int(r.Int31n(100)) + 2 // Traverse at least two points
	var points []*routeguide.Point
	for i := 0; i < pointCount; i++ {
		points = append(points, randomPoint(r))
	}
	log.Printf("Traversing %d points.", len(points))
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	stream, err := c.client.RecordRoute(ctx)
	if err != nil {
		log.Fatalf("%v.RecordRoute(_) = _, %v", c.client, err)
	}
	for _, point := range points {
		if err := stream.Send(point); err != nil {
			log.Fatalf("%v.Send(%v) = %v", stream, point, err)
		}
	}
	reply, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatalf("%v.CloseAndRecv() got error %v, want %v", stream, err, nil)
	}
	log.Printf("Route summary: %v", reply)
}

func (c *Client) RunRouteChat() {
	notes := []*routeguide.RouteNote{
		{Location: &routeguide.Point{Latitude: 0, Longitude: 1}, Message: "First message"},
		{Location: &routeguide.Point{Latitude: 0, Longitude: 2}, Message: "Second message"},
		{Location: &routeguide.Point{Latitude: 0, Longitude: 3}, Message: "Third message"},
		{Location: &routeguide.Point{Latitude: 0, Longitude: 1}, Message: "Fourth message"},
		{Location: &routeguide.Point{Latitude: 0, Longitude: 2}, Message: "Fifth message"},
		{Location: &routeguide.Point{Latitude: 0, Longitude: 3}, Message: "Sixth message"},
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	stream, err := c.client.RouteChat(ctx)
	if err != nil {
		log.Fatalf("%v.RouteChat(_) = _, %v", c.client, err)
	}
	waitc := make(chan struct{})
	go func() {
		for {
			in, err := stream.Recv()
			if err == io.EOF {
				// read done.
				close(waitc)
				return
			}
			if err != nil {
				log.Fatalf("Failed to receive a note : %v", err)
			}
			log.Printf("Got message %s at point(%d, %d)", in.Message, in.Location.Latitude, in.Location.Longitude)
		}
	}()
	for _, note := range notes {
		if err := stream.Send(note); err != nil {
			log.Fatalf("Failed to send a note: %v", err)
		}
	}
	stream.CloseSend()
	<-waitc
}

func randomPoint(r *rand.Rand) *routeguide.Point {
	lat := (r.Int31n(180) - 90) * 1e7
	long := (r.Int31n(360) - 180) * 1e7
	return &routeguide.Point{Latitude: lat, Longitude: long}
}
