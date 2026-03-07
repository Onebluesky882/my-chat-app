package db

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/gocql/gocql"
	"github.com/joho/godotenv"
)

func ConnectScylla() (*gocql.Session, error) {
	// cluster := gocql.NewCluster(
	// 	"node-0.aws-ap-southeast-1.fbc03f48f0d38c529b59.clusters.scylla.cloud:9042",
	// 	"node-1.aws-ap-southeast-1.fbc03f48f0d38c529b59.clusters.scylla.cloud:9042",
	// 	"node-2.aws-ap-southeast-1.fbc03f48f0d38c529b59.clusters.scylla.cloud:9042",
	// )
	// cluster.Authenticator = gocql.PasswordAuthenticator{Username: "scylla", Password: "EwmDNfvzMde158n"}
	// cluster.PoolConfig.HostSelectionPolicy = gocql.TokenAwareHostPolicy(gocql.DCAwareRoundRobinPolicy("AWS_AP_SOUTHEAST_1"))

	_ = godotenv.Load()

	hosts := os.Getenv("SCYLLA_HOSTS")
	username := os.Getenv("SCYLLA_USERNAME")
	password := os.Getenv("SCYLLA_PASSWORD")
	keyspace := os.Getenv("SCYLLA_KEYSPACE")
	dc := os.Getenv("SCYLLA_DATACENTER")

	if hosts == "" {
		return nil, fmt.Errorf("SCYLLA_HOSTS not set")
	}
	hostList := strings.Split(hosts, ",")

	cluster := gocql.NewCluster(hostList...)
	cluster.Port = 9042
	cluster.Keyspace = keyspace

	cluster.Authenticator = gocql.PasswordAuthenticator{
		Username: username,
		Password: password,
	}
	cluster.Timeout = 10 * time.Second
	cluster.ConnectTimeout = 10 * time.Second

	cluster.Consistency = gocql.LocalQuorum
	cluster.DisableInitialHostLookup = true

	cluster.PoolConfig.HostSelectionPolicy =
		gocql.TokenAwareHostPolicy(
			gocql.DCAwareRoundRobinPolicy(dc),
		)

	session, err := cluster.CreateSession()
	if err != nil {
		return nil, fmt.Errorf("scylla connect error: %w", err)
	}

	fmt.Println("✅ Connected to ScyllaDB")

	return session, nil
}
