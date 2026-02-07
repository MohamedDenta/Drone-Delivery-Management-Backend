package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Client struct {
	rdb *redis.Client
}

func NewClient(addr string, password string, db int) (*Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password, // no password set
		DB:       db,       // use default DB
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to redis: %w", err)
	}

	return &Client{rdb: rdb}, nil
}

func (c *Client) SetDroneLocation(ctx context.Context, id string, lat, lon float64) error {
	key := fmt.Sprintf("drone:%s:location", id)
	val := fmt.Sprintf("%f,%f", lat, lon)
	return c.rdb.Set(ctx, key, val, 1*time.Minute).Err() // TTL 1 minute
}

func (c *Client) SetDroneHeartbeat(ctx context.Context, id string) error {
	key := fmt.Sprintf("drone:%s:heartbeat", id)
	return c.rdb.Set(ctx, key, "alive", 30*time.Second).Err() // 30s TTL
}

func (c *Client) HasDroneHeartbeat(ctx context.Context, id string) (bool, error) {
	key := fmt.Sprintf("drone:%s:heartbeat", id)
	exists, err := c.rdb.Exists(ctx, key).Result()
	return exists > 0, err
}

func (c *Client) GetDroneLocation(ctx context.Context, id string) (float64, float64, error) {
	key := fmt.Sprintf("drone:%s:location", id)
	val, err := c.rdb.Get(ctx, key).Result()
	if err != nil {
		return 0, 0, err
	}

	var lat, lon float64
	_, err = fmt.Sscanf(val, "%f,%f", &lat, &lon)
	if err != nil {
		return 0, 0, fmt.Errorf("failed to parse location: %w", err)
	}
	return lat, lon, nil
}

func (c *Client) Close() error {
	return c.rdb.Close()
}
