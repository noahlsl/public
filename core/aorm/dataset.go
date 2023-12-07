package aorm

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/doug-martin/goqu/v9"

	_ "github.com/doug-martin/goqu/v9/dialect/mysql"
	"github.com/doug-martin/goqu/v9/exp"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
	"github.com/zeromicro/go-zero/core/logx"
	"gitlab.galaxy123.cloud/base/public/helper/slicex"
	"gitlab.galaxy123.cloud/base/public/helper/sqx"
)

type Dataset[T BaseModel] struct {
	client        *Client
	selectDataset *goqu.SelectDataset
	deleteDataset *goqu.DeleteDataset
	insertDataset *goqu.InsertDataset
	updateDataset *goqu.UpdateDataset
	expressions   []goqu.Expression
	having        []goqu.Expression
	rows          []interface{}
	setValues     interface{}
	asc           []string
	desc          []string
	groupBy       []string
	page          uint64
	limit         uint64
	tableName     string
	mod           T
}

func NewDataset[T BaseModel](client *Client) *Dataset[T] {
	return &Dataset[T]{
		client: client,
	}
}

func (s *Dataset[T]) toSql() (string, error) {
	var (
		sqlStr string
		err    = errors.New("no dataset")
	)

	if s.selectDataset != nil {
		se := s.selectDataset.Where(s.expressions...)
		if s.tableName != "" {
			se = se.From(s.tableName)
		}
		if s.limit != 0 {
			se = se.Limit(uint(s.limit))
		}
		if s.page != 0 {
			if s.limit == 0 {
				s.limit = 20
			}
			s.page -= 1
			se = se.Offset(uint(s.page * s.limit))
		}
		for _, s2 := range s.asc {
			se = se.Order(goqu.C(s2).Asc())
		}
		for _, s2 := range s.desc {
			se = se.Order(goqu.C(s2).Desc())
		}
		for _, s2 := range s.groupBy {
			se = se.GroupBy(s2)
		}
		for _, s2 := range s.having {
			se = se.Having(s2)
		}
		sqlStr, _, err = se.ToSQL()

		// 删除语句
	} else if s.deleteDataset != nil {
		if len(s.expressions) == 0 {
			panic("no condition in delete statement")
		}
		se := s.deleteDataset.Where(s.expressions...)
		if s.tableName != "" {
			se = se.From(s.tableName)
		}
		if s.limit != 0 {
			se = se.Limit(uint(s.limit))
		}
		sqlStr, _, err = se.ToSQL()

		// 更新语句
	} else if s.updateDataset != nil {
		if len(s.expressions) == 0 {
			panic("no condition in update statement")
		}
		se := s.updateDataset.Where(s.expressions...)
		se = se.Set(s.setValues)
		if s.tableName != "" {
			se = se.From(s.tableName)
		}
		if s.limit != 0 {
			se = se.Limit(uint(s.limit))
		}
		sqlStr, _, err = se.ToSQL()

		// 插入语句
	} else if s.insertDataset != nil {
		se := s.insertDataset
		if s.tableName != "" {
			se = s.client.g.From(s.tableName).Insert()
		}
		se = se.Rows(s.rows)
		sqlStr, _, err = se.ToSQL()
	}
	logx.WithCallerSkip(2).Debug(sqlStr)
	return sqlStr, err
}
func (s *Dataset[T]) clean() *Dataset[T] {
	s.selectDataset = nil
	s.deleteDataset = nil
	s.insertDataset = nil
	s.updateDataset = nil
	s.expressions = nil
	s.having = nil
	s.rows = nil
	s.setValues = nil
	s.asc = nil
	s.desc = nil
	s.groupBy = nil
	s.page = 0
	s.limit = 0
	s.tableName = ""
	return s
}
func (s *Dataset[T]) TableName(name string) {
	s.tableName = name
}

func (s *Dataset[T]) SQL() string {
	toSql, err := s.toSql()
	if err != nil {
		panic(err)
	}
	return toSql
}
func (s *Dataset[T]) Where(expressions ...goqu.Expression) *Dataset[T] {
	s.expressions = append(s.expressions, expressions...)
	return s
}
func (s *Dataset[T]) CopyWhere() []goqu.Expression {
	return s.expressions
}
func (s *Dataset[T]) In(field string, params ...interface{}) *Dataset[T] {
	if len(params) == 0 {
		return s
	}

	var val []interface{}
	for _, param := range params {
		ps, ok := slicex.CreateAnyTypeSlice(param)
		if ok {
			val = append(val, ps...)
			continue
		}
		val = append(val, param)
	}

	s.expressions = append(s.expressions, goqu.C(field).In(val...))
	return s
}
func (s *Dataset[T]) NotIn(field string, params interface{}) *Dataset[T] {
	s.expressions = append(s.expressions, goqu.C(field).NotIn(params))
	return s
}
func (s *Dataset[T]) Gt(field string, params interface{}) *Dataset[T] {
	s.expressions = append(s.expressions, goqu.C(field).Gt(params))
	return s
}
func (s *Dataset[T]) Gte(field string, params interface{}) *Dataset[T] {
	s.expressions = append(s.expressions, goqu.C(field).Gte(params))
	return s
}
func (s *Dataset[T]) Lt(field string, params interface{}) *Dataset[T] {
	s.expressions = append(s.expressions, goqu.C(field).Lt(params))
	return s
}
func (s *Dataset[T]) Lte(field string, params interface{}) *Dataset[T] {
	s.expressions = append(s.expressions, goqu.C(field).Lte(params))
	return s
}
func (s *Dataset[T]) Eq(field string, params interface{}) *Dataset[T] {
	s.expressions = append(s.expressions, goqu.C(field).Eq(params))
	return s
}
func (s *Dataset[T]) Neq(field string, params interface{}) *Dataset[T] {
	s.expressions = append(s.expressions, goqu.C(field).Neq(params))
	return s
}
func (s *Dataset[T]) Between(field string, start, stop interface{}) *Dataset[T] {
	s.expressions = append(s.expressions, goqu.C(field).Between(goqu.Range(start, stop)))
	return s
}
func (s *Dataset[T]) NotBetween(field string, start, stop interface{}) *Dataset[T] {
	s.expressions = append(s.expressions, goqu.C(field).NotBetween(goqu.Range(start, stop)))
	return s
}

// Like 两边模糊匹配
func (s *Dataset[T]) Like(field string, param interface{}) *Dataset[T] {
	s.expressions = append(s.expressions, goqu.C(field).Like("%"+fmt.Sprintf("%v", param)+"%"))
	return s
}

// LLike 左边模糊匹配
func (s *Dataset[T]) LLike(field string, param interface{}) *Dataset[T] {
	s.expressions = append(s.expressions, goqu.C(field).Like("%"+fmt.Sprintf("%v", param)))
	return s
}

// RLike 右边模糊匹配
func (s *Dataset[T]) RLike(field string, param interface{}) *Dataset[T] {
	s.expressions = append(s.expressions, goqu.C(field).Like(fmt.Sprintf("%v", param)+"%"))
	return s
}
func (s *Dataset[T]) Join(table exp.Expression, condition exp.JoinCondition) *Dataset[T] {
	s.InnerJoin(table, condition)
	return s
}

func (s *Dataset[T]) InnerJoin(table exp.Expression, condition exp.JoinCondition) *Dataset[T] {
	s.selectDataset.InnerJoin(table, condition)
	return s
}

func (s *Dataset[T]) FullOuterJoin(table exp.Expression, condition exp.JoinCondition) *Dataset[T] {
	s.selectDataset.FullOuterJoin(table, condition)
	return s
}

func (s *Dataset[T]) RightOuterJoin(table exp.Expression, condition exp.JoinCondition) *Dataset[T] {
	s.selectDataset.RightOuterJoin(table, condition)
	return s
}

func (s *Dataset[T]) LeftOuterJoin(table exp.Expression, condition exp.JoinCondition) *Dataset[T] {
	s.selectDataset.LeftOuterJoin(table, condition)
	return s
}

func (s *Dataset[T]) FullJoin(table exp.Expression, condition exp.JoinCondition) *Dataset[T] {
	s.selectDataset.FullJoin(table, condition)
	return s
}

func (s *Dataset[T]) RightJoin(table exp.Expression, condition exp.JoinCondition) *Dataset[T] {
	s.selectDataset.RightJoin(table, condition)
	return s
}

func (s *Dataset[T]) LeftJoin(table exp.Expression, condition exp.JoinCondition) *Dataset[T] {
	s.selectDataset.LeftJoin(table, condition)
	return s
}

func (s *Dataset[T]) NaturalJoin(table exp.Expression) *Dataset[T] {
	s.selectDataset.NaturalJoin(table)
	return s
}

func (s *Dataset[T]) NaturalLeftJoin(table exp.Expression) *Dataset[T] {
	s.selectDataset.NaturalLeftJoin(table)
	return s
}

func (s *Dataset[T]) NaturalRightJoin(table exp.Expression) *Dataset[T] {
	s.selectDataset.NaturalRightJoin(table)
	return s
}

func (s *Dataset[T]) NaturalFullJoin(table exp.Expression) *Dataset[T] {
	s.selectDataset.NaturalFullJoin(table)
	return s
}

func (s *Dataset[T]) CrossJoin(table exp.Expression) *Dataset[T] {
	s.selectDataset.CrossJoin(table)
	return s
}

func (s *Dataset[T]) Limit(l uint64) *Dataset[T] {
	s.limit = l
	return s
}
func (s *Dataset[T]) Page(l uint64) *Dataset[T] {
	s.page = l
	return s
}
func (s *Dataset[T]) ASC(fields ...string) *Dataset[T] {
	s.asc = append(s.asc, fields...)
	return s
}
func (s *Dataset[T]) Desc(fields ...string) *Dataset[T] {
	s.desc = append(s.desc, fields...)
	return s
}

// GroupBy 一般配合Count或者Sum使用
func (s *Dataset[T]) GroupBy(fields ...string) *Dataset[T] {
	s.groupBy = append(s.groupBy, fields...)
	return s
}
func (s *Dataset[T]) Having(expressions ...goqu.Expression) *Dataset[T] {
	s.having = append(s.having, expressions...)
	return s
}

func (s *Dataset[T]) Exec(ctx context.Context) (sql.Result, error) {
	sqlStr, err := s.toSql()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return s.client.db.ExecContext(ctx, sqlStr)
}
func (s *Dataset[T]) TxExec(ctx context.Context, tx *sql.Tx) (sql.Result, error) {
	sqlStr, err := s.toSql()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return tx.ExecContext(ctx, sqlStr)
}
func (s *Dataset[T]) Get(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	query, err := s.toSql()
	if err != nil {
		return errors.WithStack(err)
	}

	return s.client.db.GetContext(ctx, dest, query, args...)
}

func (s *Dataset[T]) Query(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	query, err := s.toSql()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return s.client.db.QueryContext(ctx, query, args...)
}
func (s *Dataset[T]) QueryRow(ctx context.Context, query string, args ...any) *sql.Row {
	query, err := s.toSql()
	if err != nil {
		return nil
	}

	return s.client.db.QueryRowContext(ctx, query, args...)
}

func (s *Dataset[T]) Queryx(ctx context.Context, query string, args ...any) (*sqlx.Rows, error) {
	query, err := s.toSql()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return s.client.db.QueryxContext(ctx, query, args...)
}

func (s *Dataset[T]) QueryxRow(ctx context.Context, query string, args ...any) *sqlx.Row {
	query, err := s.toSql()
	if err != nil {
		return nil
	}

	return s.client.db.QueryRowxContext(ctx, query, args...)
}
func (s *Dataset[T]) Select(fields ...interface{}) *Dataset[T] {
	s.clean()
	if len(fields) == 0 {
		s.selectDataset = s.client.g.From(s.mod.TableName()).Select(sqx.GetFields(s.mod)...)
	} else {
		s.selectDataset = s.client.g.From(s.mod.TableName()).Select(fields...)
	}

	return s
}
func (s *Dataset[T]) Update(in interface{}) *Dataset[T] {
	s.clean()
	s.setValues = in
	s.updateDataset = s.client.g.From(s.mod.TableName()).Update()
	return s
}
func (s *Dataset[T]) Install(in interface{}) *Dataset[T] {
	s.clean()
	s.rows = append(s.rows, in)
	s.insertDataset = s.client.g.From(s.mod.TableName()).Insert()
	return s
}
func (s *Dataset[T]) Delete() *Dataset[T] {
	s.clean()
	s.deleteDataset = s.client.g.From(s.mod.TableName()).Delete()
	return s
}
func (s *Dataset[T]) FindOne(ctx context.Context) (T, error) {
	var out T
	query, err := s.toSql()
	if err != nil {
		return out, errors.WithStack(err)
	}

	err = s.client.db.GetContext(ctx, &out, query)
	if err != nil {
		return out, errors.WithStack(err)
	}

	return out, nil
}
func (s *Dataset[T]) FindAll(ctx context.Context) ([]T, error) {

	query, err := s.toSql()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	var out []T
	err = s.client.db.SelectContext(ctx, &out, query)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return out, nil
}
func (s *Dataset[T]) Count(ctx context.Context, fields ...string) (int64, error) {
	s.clean()
	var field = "1"
	if len(fields) != 0 {
		field = fields[0]
	}
	s.selectDataset = s.client.g.From(s.mod.TableName()).Select(goqu.COUNT(field))
	query, err := s.toSql()
	if err != nil {
		return 0, errors.WithStack(err)
	}

	var out int64
	err = s.client.db.GetContext(ctx, &out, query)
	if err != nil {
		return 0, errors.WithStack(err)
	}

	return out, nil
}
func (s *Dataset[T]) Sum(ctx context.Context, field string) (float64, error) {

	s.clean()
	s.selectDataset = s.client.g.From(s.mod.TableName()).Select(goqu.SUM(field))
	query, err := s.toSql()
	if err != nil {
		return 0, errors.WithStack(err)
	}

	var out float64
	err = s.client.db.GetContext(ctx, &out, query)
	if err != nil {
		return 0, errors.WithStack(err)
	}

	return out, nil
}
