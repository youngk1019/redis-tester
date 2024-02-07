package internal

import (
	"fmt"
	"net"
	"strings"

	testerutils "github.com/codecrafters-io/tester-utils"
)

func testReplInfoReplica(stageHarness *testerutils.StageHarness) error {
	listener, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println("Error starting TCP server:", err)
	}
	defer listener.Close()
	logger := stageHarness.Logger

	logger.Infof("Server is running on port 6379")

	replica := NewRedisBinary(stageHarness)
	replica.args = []string{
		"--port", "6380",
		"--replicaof", "localhost", "6379",
	}

	if err := replica.Run(); err != nil {
		return err
	}

	conn, err := listener.Accept()
	if err != nil {
		fmt.Println("Error accepting: ", err.Error())
		return err
	}
	client := NewRedisClient("localhost:6380")

	logger.Infof("$ redis-cli INFO replication")
	resp, err := client.Info("replication").Result()
	lines := strings.Split(resp, "\n")
	infoMap := parseInfoOutput(lines, ":")
	key := "role"
	role := infoMap[key]

	if err != nil {
		logFriendlyError(logger, err)
		return err
	}

	if key != "role" {
		return fmt.Errorf("Expected 'role' key in INFO replication, got %v", key)
	}

	if role != "slave" {
		return fmt.Errorf("Expected 'role' to be 'slave' in INFO replication, got %v", role)
	}

	client.Close()
	conn.Close()
	return nil
}
