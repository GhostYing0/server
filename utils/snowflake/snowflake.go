package snowflake

import (
	"github.com/bwmarrin/snowflake"
	"time"
)

// GenSnowflakeID 生成雪花算法ID
func GenSnowflakeID() int64 {
	var node *snowflake.Node
	var st time.Time
	st, _ = time.Parse("2006-01-02", "2023-03-22")
	snowflake.Epoch = st.UnixNano() / 1000000
	node, _ = snowflake.NewNode(123)
	return node.Generate().Int64() % 1e10
}
