package kfkmodule

const topicPrefix = "user_wallet_result_topic_"

func GetWalletChangeTopic(suffix string) string {
	if len(suffix) == 0 {
		panic("GetWalletChangeTopic: empty suffix")
	}
	return topicPrefix + suffix
}

type PushData struct {
	UID     uint32
	Balance int64
}
