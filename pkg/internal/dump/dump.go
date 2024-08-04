package dump

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	atomicfw "github.com/natefinch/atomic"
	"os"
)

type UserWalletResultDump struct {
	Res       map[uint32]UserWalletResult
	KfkOffset int64
}

type UserWalletResult struct {
	UID     uint32
	Balance int64
	Change  bool `json:"-"`
}

const dumpDir = "datadir/user_wallet_dump/"

func LoadDumpData(sharding int32) UserWalletResultDump {
	dump := UserWalletResultDump{}
	tmpData, err := os.ReadFile(GetDumpFile(sharding))
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(tmpData, &dump)
	if err != nil {
		panic(err)
	}
	if dump.KfkOffset == 0 {
		dump.KfkOffset = sarama.OffsetNewest
	}
	return dump
}

func GetDumpFile(sharding int32) string {
	return fmt.Sprint(dumpDir, "user_sharding_", sharding, ".json")
}

func MakeDumpWithUserSharding(res UserWalletResultDump, sharding int32) error {
	_, err := os.Stat(dumpDir)
	if os.IsNotExist(err) {
		if err := os.MkdirAll(dumpDir, os.ModePerm); err != nil {
			return err
		}
	}

	dumpData, _ := json.Marshal(res)
	filePath := GetDumpFile(sharding)
	if err = atomicfw.WriteFile(filePath, bytes.NewBuffer(dumpData)); err != nil {
		return err
	}

	// -rw-r--r--
	os.Chmod(filePath, 0644)
	return nil
}
