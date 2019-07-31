package db

import (
	"testing"
	"time"
)

// VO value object
type VO struct {
	BaseVO
	Code       string    `json:"code"`
	Name       string    `json:"name"`
	FinishedAt time.Time `json:"finished_at"`
}

func TestVO(t *testing.T) {
	now := time.Now()
	row := MapRow{"code": "code1", "name": "tom01", "finished_at": now}
	vo := VO{}
	row.CopyToVO(&vo)
	if vo.Code != row.String("code") && vo.FinishedAt.String() != now.String() {
		t.Fatal()
	}

	vo = VO{}
	row.CopyToWithJSON(&vo)
	if vo.Code != row.String("code") && vo.FinishedAt.String() != now.String() {
		t.Fatal()
	}
}

func TestIsEmpty(t *testing.T) {
	vo := VO{}
	cv := NewStructConvert(vo)
	row := cv.Struct2DataRow(ActionInsert)

	if row.IsEmpty() {
		t.Fatal("is empty failed", row)
	}

	vo.Name = "tom"
	vo.FinishedAt = time.Now()
	cv = NewStructConvert(vo)
	row = cv.Struct2DataRow(ActionUpdate)

	if len(row) != 3 {
		t.Fatal("is empty failed", row)
	}

}
