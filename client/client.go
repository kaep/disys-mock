package main

import (
	"context"
	"fmt"
	"log"
	m "mock/proto"
	"time"

	"google.golang.org/grpc"
)

type Frontend struct {
	ctx       context.Context
	timestamp int32
	servers   map[int]m.MockClient
}

type Client struct {
	Id    int32
	Front Frontend
}

func main() {
	log.Print("Starting client")
	client := Client{}
	client.Id = 1 //husk at der skal laves flere klienter med forskellige id?

	client.setupFrontend()
	log.Printf("Frontend has been set up on client %v", client.Id)

	client.Increment()
	time.Sleep(2 * time.Second)
	client.Increment()
	time.Sleep(2 * time.Second)
	client.Increment()
	time.Sleep(2 * time.Second)
}

func (c *Client) Increment() {
	c.Front.timestamp++
	for i := range c.Front.servers {
		value, err := c.Front.servers[i].Increment(c.Front.ctx, &m.Request{Id: c.Id})
		if err != nil {
			log.Printf("Error attempting to reach server %v - removing", i)
			log.Printf("Error: %v", err)
			delete(c.Front.servers, i)
		} else {
			log.Printf("Client %v recieved value %v from server %v", c.Id, value.Value, i)
		}
	}
}

func (c *Client) setupFrontend() {

	//HUSK AT CONNECTION SKAL DEFER CLOSE()

	var ports = 8080
	var numServers = 1 //HUSK AT ÆNDRE TIL FLERE !
	c.Front.servers = make(map[int]m.MockClient)
	for i := 0; i < numServers; i++ {
		var conn *grpc.ClientConn
		var port = fmt.Sprintf(":%v", ports+i)
		conn, err := grpc.Dial(port, grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			log.Fatalf("Could not connect: %s", err)
		}
		client := m.NewMockClient(conn)
		c.Front.servers[i] = client
	}
	ctx := context.Background()
	c.Front.ctx = ctx
}
