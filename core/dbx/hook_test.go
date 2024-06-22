package dbx

import (
	"fmt"
	"testing"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// Bank 银行表
type Bank struct {
	ID        int64  `gorm:"column:id;type:bigint;primaryKey;autoIncrement:true;comment:系统银行id" json:"id"`                             // 银行id
	UpdatedAt int64  `gorm:"column:updated_at;type:bigint;type:unsigned;not null;autoUpdateTime:milli;comment:更新时间" json:"updated_at"` // 更新时间
	Operator  string `gorm:"column:operator;type:varchar(64);not null;comment:操作人" json:"operator"`                                    // 操作人
}

// TableName Bank's table name
func (*Bank) TableName() string {
	return "bank"
}

func TestHook(t *testing.T) {
	dsn := "mall:kdi8WMe8HfXfZS5N@tcp(127.0.0.1:3306)/mall?charset=utf8mb4&parseTime=true&loc=Asia%2FShanghai"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Log(err)
		return
	}
	// 注册插件
	err = db.Use(&ReplaceSelectStatementPlugin{})
	if err != nil {
		t.Log(err)
		return
	}
	var banks []Bank
	err = db.Model(&Bank{}).Select("id").Limit(1).Find(&banks).Debug().Error
	if err != nil {
		t.Log(err)
		return
	}

	//var banks []Bank
	var count int64
	err = db.Model(&Bank{}).Where("id > ?", 0).Count(&count).Debug().Error
	if err != nil {
		t.Log(err)
		return
	}

	fmt.Println("通过")
}
