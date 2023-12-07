package orm

import (
	"context"
	"fmt"
	"github.com/doug-martin/goqu/v9"
	"testing"

	"github.com/noahlsl/public/core/member"
)

func TestNewClient(t *testing.T) {
	c := NewClient("root:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=true&loc=Asia%2FShanghai")
	ex := goqu.Ex{"member_id": 113}
	ctx := context.Background()
	// 查询
	one, err := c.Member.Select("member_id", "member_name").Where(ex).FindOne(ctx)
	if err != nil {
		fmt.Println(one)
		return
	}

	fmt.Println("查询成功")

	// 更新
	_, err = c.Member.Update(map[string]interface{}{"member_name": "abe1013"}).Where(ex).Exec(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("更新成功")

	// 新增
	m := member.ModelMember{
		MemberId: 119,
	}
	_, err = c.Member.Install(m).Exec(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("新增成功")

	// 删除
	ex = goqu.Ex{"member_id": 117}
	_, err = c.Member.Delete().Where(ex).Exec(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("删除成功")
}

func TestTx(t *testing.T) {
	ex := goqu.Ex{"member_id": 120}
	ctx := context.Background()
	c := NewClient("root:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=true&loc=Asia%2FShanghai")

	tx, err := c.Tx(ctx)
	if err != nil {
		fmt.Println(err)
		return
	}

	// 更新
	tx.Add(c.Member.Update(map[string]interface{}{"member_id": "1"}).Where(ex))

	// 新增
	m := member.ModelMember{
		MemberId: 121,
	}
	tx.Add(c.Member.Install(m))
	// 删除
	tx.Add(c.Member.Delete().Where(ex))
	err = tx.Commit()
	if err != nil {
		fmt.Println(err)
		return
	}
	_ = c.Close()
}
