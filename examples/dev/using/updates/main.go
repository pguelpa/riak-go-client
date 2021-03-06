package main

import (
	"fmt"

	riak "github.com/basho/riak-go-client"
)

/*
    Code samples from:
    http://docs.basho.com/riak/latest/dev/using/updates/

	make sure this bucket-type is created:
	siblings

	riak-admin bucket-type create siblings '{"props":{"allow_mult":true}}'
	riak-admin bucket-type activate siblings
*/
func main() {
	//riak.EnableDebugLogging = true

	nodeOpts := &riak.NodeOptions{
		RemoteAddress: "riak-test:10017",
	}

	var node *riak.Node
	var err error
	if node, err = riak.NewNode(nodeOpts); err != nil {
		fmt.Println(err.Error())
	}

	nodes := []*riak.Node{node}
	opts := &riak.ClusterOptions{
		Nodes: nodes,
	}

	cluster, err := riak.NewCluster(opts)
	if err != nil {
		fmt.Println(err.Error())
	}

	defer func() {
		if err := cluster.Stop(); err != nil {
			fmt.Println(err.Error())
		}
	}()

	if err := cluster.Start(); err != nil {
		fmt.Println(err.Error())
	}

	// ping
	ping := &riak.PingCommand{}
	if err := cluster.Execute(ping); err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("ping passed")
	}

	storeCoach(cluster)

	if err := updateCoach(cluster, "seahawks", "Bob Abooey"); err != nil {
		fmt.Println(err.Error())
	}
}

func storeCoach(cluster *riak.Cluster) {
	obj := &riak.Object{
		ContentType:     "text/plain",
		Charset:         "utf-8",
		ContentEncoding: "utf-8",
		Value:           []byte("Pete Carroll"),
	}

	cmd, err := riak.NewStoreValueCommandBuilder().
		WithBucketType("siblings").
		WithBucket("coaches").
		WithKey("seahawks").
		WithContent(obj).
		Build()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	if err := cluster.Execute(cmd); err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("Stored Pete Carroll")
}

func updateCoach(cluster *riak.Cluster, team, newCoach string) error {
	var cmd riak.Command
	var err error

	cmd, err = riak.NewFetchValueCommandBuilder().
		WithBucketType("siblings").
		WithBucket("coaches").
		WithKey(team).
		Build()

	if err != nil {
		return err
	}

	if err := cluster.Execute(cmd); err != nil {
		return err
	}

	fvc := cmd.(*riak.FetchValueCommand)
	obj := fvc.Response.Values[0]
	obj.Value = []byte(newCoach)

	cmd, err = riak.NewStoreValueCommandBuilder().
		WithBucketType("siblings").
		WithBucket("coaches").
		WithKey(team).
		WithContent(obj).
		Build()

	if err != nil {
		return err
	}

	if err := cluster.Execute(cmd); err != nil {
		return err
	}

	return nil
}
