package wb

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"gorm.io/gorm"
	"log"
	"sync"
	"time"
	"user_wallet/pkg/internal/dump"
	"user_wallet/pkg/internal/kfkmodule"
	"user_wallet/pkg/internal/model"
)

type WbService struct {
	dumpCh       chan dump.UserWalletResultDump
	userSharding int32

	cacheMux              sync.RWMutex
	userWalletResultCache UserWalletResultCache // userid => balance info
	lastKfakOffset        int64
}

func NewWbService(userSharding int32) *WbService {
	d := dump.LoadDumpData(userSharding)
	ws := &WbService{
		dumpCh:                make(chan dump.UserWalletResultDump, 8),
		userWalletResultCache: make(UserWalletResultCache, len(d.Res)),
		userSharding:          userSharding,
	}

	for _, result := range d.Res {
		fmt.Println("result:", result.UID)
		ws.userWalletResultCache[result.UID] = &dump.UserWalletResult{
			UID:     result.UID,
			Balance: result.Balance,
		}
	}

	log.Println("ws.userWalletResultCache--", len(ws.userWalletResultCache))

	ws.lastKfakOffset = d.KfkOffset

	go ws.LoopDump()
	go ws.LoopWriteDB()
	return ws
}

func (ws *WbService) PutCache(uid uint32, balance int64, offset int64) {
	ws.cacheMux.Lock()
	if old, ok := ws.userWalletResultCache[uid]; ok {
		if old.Balance != balance {
			old.Balance = balance
			old.Change = true
		}
	}
	ws.lastKfakOffset = offset
	ws.cacheMux.Unlock()
}

func (ws *WbService) WriteDBTx() {
	log.Println("start write db transaction")
	ws.cacheMux.Lock()
	defer ws.cacheMux.Unlock()

	if len(ws.userWalletResultCache) == 0 {
		return
	}
	dumpData := dump.UserWalletResultDump{Res: map[uint32]dump.UserWalletResult{}}

	dumpData.KfkOffset = ws.lastKfakOffset

	err := model.GetDB().Transaction(func(tx *gorm.DB) error {
		for _, rd := range ws.userWalletResultCache {
			dumpData.Res[rd.UID] = dump.UserWalletResult{UID: rd.UID, Balance: rd.Balance}
			if !rd.Change {
				continue
			}
			if err1 := model.UpdateBalance(rd.UID, rd.Balance, tx); err1 != nil {
				return err1
			}
			rd.Change = false
		}
		return nil
	})

	if len(dumpData.Res) == 0 {
		return
	}

	if err != nil {
		panic(err)
	}

	log.Printf("dumpData: %d, kfkOffset: %d", len(dumpData.Res), dumpData.KfkOffset)

	ws.dumpCh <- dumpData
}

func (ws *WbService) LoopWriteDB() {
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()
	for range ticker.C {
		ws.WriteDBTx()
	}
}

func (ws *WbService) LoopDump() {
	lastDumpTs := time.Now()

	for v := range ws.dumpCh {
		now := time.Now()
		if now.Sub(lastDumpTs) <= 3*time.Second {
			continue
		}
		err := dump.MakeDumpWithUserSharding(v, ws.userSharding)
		if err != nil {
			log.Printf("dumpData error: %v", err)
		}
		lastDumpTs = now
	}
}

func (ws *WbService) Consume() {
	log.Printf("Consume kafka from offset: %d", ws.lastKfakOffset)
	err := model.Init()
	if err != nil {
		log.Fatal(err)
	}

	consumer, err := kfkmodule.NewConsumer(ws.userSharding, nil, ws.lastKfakOffset)
	if err != nil {
		log.Fatal(err)
	}
	err = consumer.Start(func(ctx context.Context, message *sarama.ConsumerMessage) {
		rd := &dump.UserWalletResult{}
		_ = json.Unmarshal(message.Value, rd)
		ws.PutCache(rd.UID, rd.Balance, message.Offset)
	})
	if err != nil {
		log.Fatal(err)
	}
}
