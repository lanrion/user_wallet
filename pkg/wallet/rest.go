package wallet

import (
	"errors"
	"fmt"
	"github.com/shopspring/decimal"
	"log"
	"net/http"
	"strconv"
	"user_wallet/pkg/internal/dump"
	"user_wallet/pkg/internal/kfkmodule"
)

var (
	reqConsumer  Consumer
	userSharding int32
)

func Run(userShard int32) {
	userSharding = userShard

	d := dump.LoadDumpData(userSharding)
	walletData := make(map[uint32]*Wallet, len(d.Res))
	for _, ud := range d.Res {
		walletData[ud.UID] = &Wallet{uid: ud.UID, balance: ud.Balance}
	}

	BucketManager = NewBucket(walletData)

	kfkProducer, err := kfkmodule.NewProducer(userShard, []string{"localhost:9092"})
	if err != nil {
		log.Fatal(err)
	}
	ps := NewPusher(userShard, kfkProducer, 10000)
	reqConsumer = NewConsumer(ps)
	err = reqConsumer.Start()
	if err != nil {
		log.Fatalf("Failed to start consumer: %v", err)
	}

	sm := http.NewServeMux()
	sm.HandleFunc("/user/info/", getUserInfoHandler)

	sm.HandleFunc("/user/deposit/", depositHandler)
	sm.HandleFunc("/user/withdraw/", withdrawHandler)

	sm.HandleFunc("/user/transfer/", transferHandler)

	port := ":9102"
	log.Printf("Starting REST server on port %s  with user sharding %d \n", port, userSharding)

	err = http.ListenAndServe(port, sm)
	if err != nil {
		log.Fatalf("Failed to start REST server: %v", err)
	}

	log.Println("rest server stopped")
}

func transferHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[transferHandler]Request from %s %s \n", r.URL, r.Method)
	uid, amount, err := getUserIdAndAmount(r)
	if err != nil {
		_, _ = w.Write([]byte("err: " + err.Error()))
		return
	}

	touid, err := getToUidFromReq(r)
	if err != nil {
		_, _ = w.Write([]byte("err: " + err.Error()))
		return
	}

	if err = checkUserSharding(int32(uid)); err != nil {
		_, _ = w.Write([]byte("err: " + err.Error()))
		return
	}

	j := NewJob(uid, amount, Transfer, touid)
	err = reqConsumer.PushBack(r.Context(), j)
	if err != nil {
		_, _ = w.Write([]byte("err: " + err.Error()))
		return
	}
	j.Wait()

	if j.Err != nil {
		_, _ = w.Write([]byte("err: " + j.Err.Error()))
		return
	}

	_, _ = w.Write([]byte("Deposit successfully"))
}

func getUserInfoHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("Request from %s %s \n", r.URL, r.Method)
	uid, err := getUserIDFromReq(r)
	if err != nil {
		_, _ = w.Write([]byte("err: " + err.Error()))
		return
	}

	if err = checkUserSharding(int32(uid)); err != nil {
		_, _ = w.Write([]byte("err: " + err.Error()))
		return
	}

	data := fmt.Sprintf("uid: %d, balance: %v", uid, BucketManager.Fetch(uid).Balance())
	_, _ = w.Write([]byte(data))
}

func depositHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[DepositHandler]Request from %s %s \n", r.URL, r.Method)

	uid, amount, err := getUserIdAndAmount(r)
	if err != nil {
		_, _ = w.Write([]byte("err: " + err.Error()))
		return
	}

	if err = checkUserSharding(int32(uid)); err != nil {
		_, _ = w.Write([]byte("err: " + err.Error()))
		return
	}

	log.Printf("Request from %s %s, %v, %d \n", r.RemoteAddr, r.Method, amount, uid)

	j := NewJob(uid, amount, Deposit, 0)
	err = reqConsumer.PushBack(r.Context(), j)
	if err != nil {
		_, _ = w.Write([]byte("err: " + err.Error()))
		return
	}
	j.Wait()

	if j.Err != nil {
		_, _ = w.Write([]byte("err: " + j.Err.Error()))
		return
	}

	_, _ = w.Write([]byte("Deposit successfully"))
}

func withdrawHandler(w http.ResponseWriter, r *http.Request) {
	log.Printf("[WithdrawHandler]Request from %s %s \n", r.URL, r.Method)
	uid, amount, err := getUserIdAndAmount(r)
	if err != nil {
		_, _ = w.Write([]byte("err: " + err.Error()))
		return
	}

	if err = checkUserSharding(int32(uid)); err != nil {
		_, _ = w.Write([]byte("err: " + err.Error()))
		return
	}

	j := NewJob(uid, amount, Withdraw, 0)
	err = reqConsumer.PushBack(r.Context(), j)
	if err != nil {
		_, _ = w.Write([]byte("err: " + err.Error()))
		return
	}

	j.Wait()

	if j.Err != nil {
		_, _ = w.Write([]byte("err: " + j.Err.Error()))
		return
	}

	_, _ = w.Write([]byte("Withdraw successfully"))
}

func getUserIdAndAmount(r *http.Request) (uint32, float64, error) {
	amount, err := getAmountFromReq(r)
	if err != nil {
		return 0, 0, err
	}

	uid, err := getUserIDFromReq(r)
	if err != nil {
		return 0, 0, err
	}

	if amount <= 0 {
		return uid, 0, errors.New("amount is negative")
	}

	return uid, amount, nil
}

func getUserIDFromReq(r *http.Request) (uint32, error) {
	uidStr := r.URL.Query().Get("user_id")
	uid, err := strconv.ParseUint(uidStr, 10, 32)
	return uint32(uid), err
}

func getToUidFromReq(r *http.Request) (uint32, error) {
	uidStr := r.URL.Query().Get("to_uid")
	uid, err := strconv.ParseUint(uidStr, 10, 32)
	return uint32(uid), err
}

func getAmountFromReq(r *http.Request) (float64, error) {
	uidStr := r.URL.Query().Get("amount")
	di, err := decimal.NewFromString(uidStr)
	if err != nil {
		return 0, err
	}
	return di.InexactFloat64(), nil
}

func checkUserSharding(uid int32) error {
	if (uid % 10) != userSharding {
		return fmt.Errorf("user(uid:%d) request sharing(%d) failed", uid, userSharding)
	}
	return nil
}
