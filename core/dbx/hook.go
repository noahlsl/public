package dbx

import (
	"fmt"
	"gorm.io/gorm"
	"reflect"
)

// Options 用于传递分表信息
type Options struct {
	BaseTableName string
	HashField     string
}

var (
	suffixFn = getTableSuffix
)

func getTableSuffix(id int64) string {
	return fmt.Sprintf("%d", id%10)
}

func SetSuffixFn(fn func(id int64) string) {
	suffixFn = fn
}

func RegisterHooks(db *gorm.DB, opts ...Options) error {
	for _, opt := range opts {
		// 创建操作回调
		err := db.Callback().Create().Before("gorm:create").Register("sharding:before_create", func(db *gorm.DB) {
			setShardedTable(db, opt)
		})
		if err != nil {
			return err
		}

		// 更新操作回调
		err = db.Callback().Update().Before("gorm:update").Register("sharding:before_update", func(db *gorm.DB) {
			setShardedTable(db, opt)
		})
		if err != nil {
			return err
		}

		// 删除操作回调
		err = db.Callback().Delete().Before("gorm:delete").Register("sharding:before_delete", func(db *gorm.DB) {
			setShardedTable(db, opt)
		})
		if err != nil {
			return err
		}

		// 查询操作回调
		err = db.Callback().Query().Before("gorm:query").Register("sharding:before_query", func(db *gorm.DB) {
			setShardedTable(db, opt)
		})
		if err != nil {
			return err
		}
	}

	return nil
}

func setShardedTable(db *gorm.DB, opts Options) {
	if db.Statement.Schema != nil && db.Statement.Schema.Table == opts.BaseTableName {
		var id int64
		// 获取哈希字段的值
		if db.Statement.ReflectValue.Kind() == reflect.Ptr {
			idField := db.Statement.ReflectValue.Elem().FieldByName(opts.HashField)
			if idField.IsValid() {
				id = idField.Interface().(int64)
			}
		} else {
			idField := db.Statement.ReflectValue.FieldByName(opts.HashField)
			if idField.IsValid() {
				id = idField.Interface().(int64)
			}
		}
		db.Statement.Table = opts.BaseTableName + "_" + suffixFn(id)
	}
}

// ReplaceSelectStatementPlugin GORM 替换掉 * 的Hook方法
type ReplaceSelectStatementPlugin struct {
	gorm.Plugin
}

// Name 替换 Select 语句的回调方法
func (p *ReplaceSelectStatementPlugin) Name() string {
	return "ReplaceSelectStatementPlugin"
}

// Initialize 在执行查询之前调用的回调方法
func (p *ReplaceSelectStatementPlugin) Initialize(db *gorm.DB) (err error) {
	// 在查询前注册一个钩子
	return db.Callback().Query().Before("gorm:query").Register("replaceSelectStatement", p.replaceSelectStatement)
}

// 替换 Select 语句的具体逻辑
func (p *ReplaceSelectStatementPlugin) replaceSelectStatement(db *gorm.DB) {
	// 获取表名称
	if len(db.Statement.Selects) == 0 && len(db.Statement.Clauses) == 0 {
		db.Statement.Select(db.Statement.Schema.DBNames)
	}
}
