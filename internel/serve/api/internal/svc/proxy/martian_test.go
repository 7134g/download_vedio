package proxy

import (
	"dv/internel/serve/api/internal/db"
	"dv/internel/serve/api/internal/util/model"
	"testing"
	"time"
)

func TestMartian(t *testing.T) {
	db.InitSqlite("test.sqlite")
	taskDB = model.NewTaskModel(db.GetDB())
	go MatchInformation()
	SetServeProxyAddress("http://127.0.0.1:7890", "", "")
	OpenCert()
	if err := Martian(); err != nil {
		t.Fatal(err)
	}

}

func TestGet(t *testing.T) {
	for i := 0; i < 10; i++ {
		t.Log(getNumber())
		time.Sleep(time.Second * 3)
	}
}
