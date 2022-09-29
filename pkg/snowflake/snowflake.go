package snowflake

import (
	"github.com/spf13/viper"
	"time"

	sf "github.com/bwmarrin/snowflake"
)

var node *sf.Node

func Init() (err error) {
	var st time.Time
	st, err = time.Parse("2006-01-02", viper.GetString("app.start_time"))
	if err != nil {
		return
	}
	sf.Epoch = st.UnixNano() / 1000000
	node, err = sf.NewNode(viper.GetInt64("app.machine_id"))
	return
}
func GenID() int64 {
	return node.Generate().Int64()
}
