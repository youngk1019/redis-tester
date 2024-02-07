package internal

import (
	"fmt"

	testerutils "github.com/codecrafters-io/tester-utils"
	"github.com/smallnest/resp3"
)

func testReplMasterReplconf(stageHarness *testerutils.StageHarness) error {
	master := NewRedisBinary(stageHarness)
	master.args = []string{
		"--port", "6379",
	}

	if err := master.Run(); err != nil {
		return err
	}

	logger := stageHarness.Logger

	conn, err := NewRedisConn("", "localhost:6379")
	if err != nil {
		fmt.Println("Error connecting to TCP server:", err)
	}

	r := resp3.NewReader(conn)
	w := resp3.NewWriter(conn)

	replica := FakeRedisReplica{
		Reader: r,
		Writer: w,
		Logger: logger,
	}

	err = replica.Ping()
	if err != nil {
		return err
	}

	err = replica.ReplConfPort()
	if err != nil {
		return err
	}

	conn.Close()
	return nil
}
