package snowflake

import (
	"errors"
	sf "github.com/bwmarrin/snowflake"
	"time"
)

var (
	InvalidInitParam  = errors.New("snowflake初始化失败:无效的startTime或machineID")
	InvalidTimeFormat = errors.New("snowflake初始化失败:无效的startTime格式")
)
var node *sf.Node

// Init initializes the snowflake node with a custom epoch and machine ID.
func Init(startTime string, machineID int64) (err error) {
	// Initialize a new snowflake node with a specific machine ID
	if len(startTime) == 0 || machineID < 0 {
		return InvalidInitParam
	}
	var st time.Time
	st, err = time.Parse("2006-01-02 15:04:05", startTime)
	if err != nil {
		return InvalidTimeFormat
	}
	sf.Epoch = st.UnixNano() / 1e6 // Set the epoch to the specified start time in milliseconds
	node, err = sf.NewNode(machineID)
	return
}
func GenID() (id int64) {
	return node.Generate().Int64()
}
