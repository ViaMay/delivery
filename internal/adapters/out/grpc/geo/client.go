package geo

import (
	"context"
	"delivery/internal/core/domain/model/kernel"
	"delivery/internal/core/ports"
	"delivery/internal/generated/clients/geosrv/geopb"
	"delivery/internal/pkg/errs"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"log"
	"time"
)

var _ ports.GeoClient = &Client{}

type Client struct {
	conn     *grpc.ClientConn
	pbClient geopb.GeoClient
	timeout  time.Duration
}

func NewClient(host string) (*Client, error) {
	if host == "" {
		return nil, errs.NewValueIsRequiredError("host")
	}

	conn, err := grpc.NewClient(host, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}

	pbClient := geopb.NewGeoClient(conn)

	return &Client{
		conn:     conn,
		pbClient: pbClient,
		timeout:  5 * time.Second,
	}, nil
}

func (c *Client) Close() error {
	return c.conn.Close()
}

func (c *Client) GetGeolocation(ctx context.Context, street string) (kernel.Location, error) {
	// Формируем запрос
	req := &geopb.GetGeolocationRequest{
		Street: street,
	}

	// Делаем запрос
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()
	resp, err := c.pbClient.GetGeolocation(ctx, req)
	if err != nil {
		return kernel.Location{}, err
	}

	// Создаем и возвращаем VO Geo
	location, err := kernel.NewLocation(int(resp.Location.X), int(resp.Location.Y))
	if err != nil {
		return kernel.Location{}, err
	}
	return location, nil
}
