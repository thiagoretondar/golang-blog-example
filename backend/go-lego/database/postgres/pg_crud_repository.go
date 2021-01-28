// Package postgres contains methods to handle data in postgreSQL
package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"reflect"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/reflectx"
	"github.com/thiagoretondar/golang-blog-example/backend/go-lego/database"
)

type Pg interface {
	database.CRUDRepository
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*Tx, error)
	WithTx(tx *Tx) PgTx
	GetConn() *sqlx.DB
}

type pgRepository struct {
	table   string
	session *sqlx.DB
	mapper  *reflectx.Mapper

	// sortableColumns map[string]bool
}

// NewRepository setup a new CRUD Adapter for a specific table
func NewRepository(tableName string, session *sqlx.DB) Pg {
	return &pgRepository{
		table:   tableName,
		session: session,
		mapper:  reflectx.NewMapper("db"),
	}
}

// Insert inserts a single record
func (b *pgRepository) Insert(ctx context.Context, data interface{}, lastInsertedID interface{}) error {
	columns, values := b.ExtractColumnPairs(data)

	// Prepare query
	queryBuilder := sq.
		Insert(b.table).
		Columns(columns...).
		Values(values...).
		Suffix("returning \"id\"").
		PlaceholderFormat(sq.Dollar)

	// Build SQL Query
	query, args, err := queryBuilder.ToSql()
	if err != nil {
		return fmt.Errorf("unable to build query: %w", err)
	}

	// Prepare statement and defer it's closure
	stmt, err := b.session.PreparexContext(ctx, query)
	if err != nil {
		return err
	}
	defer func(stmt *sqlx.Stmt) {
		//b.logger.SafeClose(stmt, "Unable to close statement")
	}(stmt)

	// Do the insert query
	if lastInsertedID != nil {
		// We do QueryRowxContext because Postgres doesn'b work with lastInsertedID
		err = stmt.QueryRowxContext(ctx, args...).Scan(lastInsertedID)
	} else {
		// Here we don'b need the lastInsertedID
		_, err = stmt.ExecContext(ctx, args...)
	}

	return err
}

// FindAll returns all records from database
func (b *pgRepository) FindAll(ctx context.Context, output interface{}) error {
	//TODO implement pagination

	// Prepare query
	queryBuilder := sq.Select("*").From(b.table).PlaceholderFormat(sq.Dollar)

	// Build SQL Query
	query, _, err := queryBuilder.ToSql()
	if err != nil {
		return fmt.Errorf("unable to build query: %w", err)
	}

	// Prepare statement and defer it's closure
	stmt, err := b.session.PreparexContext(ctx, query)
	if err != nil {
		return err
	}
	defer func(stmt *sqlx.Stmt) {
		//b.logger.SafeClose(stmt, "Unable to close statement")
	}(stmt)

	err = stmt.SelectContext(ctx, output)
	if err != nil {
		return err
	}
	return nil
}

// FindOne returns only one record given the filter
func (b *pgRepository) FindOne(ctx context.Context, filter interface{}, output interface{}) error {
	columns, _ := b.ExtractColumnPairs(output)

	// Prepare query
	qb := sq.Select(columns...).
		From(b.table).
		Where(filter).
		Limit(1).
		PlaceholderFormat(sq.Dollar)

	// Build SQL Query
	query, args, err := qb.ToSql()
	if err != nil {
		return err
	}

	// Prepare statement and defer it's closure
	stmt, err := b.session.PreparexContext(ctx, query)
	if err != nil {
		return err
	}
	defer func(stmt *sqlx.Stmt) {
		//b.logger.SafeClose(stmt, "Unable to close statement")
	}(stmt)

	// Do the find
	err = b.session.QueryRowxContext(ctx, query, args...).StructScan(output)
	if err != nil {
		return err
	}

	// No errors
	return nil
}

// Find returns all records from database that match the filter
func (b *pgRepository) Find(ctx context.Context, filter interface{}, output interface{}) error {
	// Prepare query
	qb := sq.Select("*").
		From(b.table).
		Where(filter).
		PlaceholderFormat(sq.Dollar)

	// Build SQL Query
	query, args, err := qb.ToSql()
	if err != nil {
		return err
	}

	// Prepare statement and defer it's closure
	stmt, err := b.session.PreparexContext(ctx, query)
	if err != nil {
		return err
	}
	defer func(stmt *sqlx.Stmt) {
		//b.logger.SafeClose(stmt, "Unable to close statement")
	}(stmt)

	// Do the find
	err = stmt.SelectContext(ctx, output, args...)
	if err != nil {
		return err
	}

	// No errors
	return nil
}

// Update updates records matching the given filter
func (b *pgRepository) Update(ctx context.Context, set map[string]interface{}, filter interface{}) (int64, error) {
	// Prepare query
	qb := sq.Update(b.table).
		SetMap(set).
		Where(filter).
		PlaceholderFormat(sq.Dollar)

	// Build SQL Query
	query, args, err := qb.ToSql()
	if err != nil {
		return 0, err
	}

	// Prepare statement and defer it's closure
	stmt, err := b.session.PreparexContext(ctx, query)
	if err != nil {
		return 0, err
	}
	defer func(stmt *sqlx.Stmt) {
		//b.logger.SafeClose(stmt, "Unable to close statement")
	}(stmt)

	// Execute query
	result, err := stmt.ExecContext(ctx, args...)
	if err != nil {
		return 0, err
	}

	// Return result and possible error
	return result.RowsAffected()
}

// Remove updates the records that match the given filter to deleted status. The record isn't really deleted from database
func (b *pgRepository) Remove(ctx context.Context, filter interface{}, physicalDeletion bool) (int64, error) {
	if !physicalDeletion {
		// Logical deletion set
		updateStatusSet := map[string]interface{}{
			// TODO remove hard coded column names
			"updated_at": "now()",
			"status":     false,
		}

		return b.Update(ctx, updateStatusSet, filter)
	}

	// Prepare query
	qb := sq.Delete(b.table).Where(filter).PlaceholderFormat(sq.Dollar)

	// Build SQL query
	query, args, err := qb.ToSql()
	if err != nil {
		return 0, err
	}

	// Prepare statement and defer it's closure
	stmt, err := b.session.PreparexContext(ctx, query)
	if err != nil {
		return 0, err
	}
	defer func(stmt *sqlx.Stmt) {
		//b.logger.SafeClose(stmt, "Unable to close statement")
	}(stmt)

	// Execute query
	result, err := stmt.ExecContext(ctx, args...)
	if err != nil {
		return 0, err
	}

	// Return result and possible error
	return result.RowsAffected()
}

// Count counts how many records match the filter. If no filter is given will return the quantity
// of all records stored
func (b *pgRepository) Count(ctx context.Context, filter interface{}) (int64, error) {
	// Prepare query
	qb := sq.Select("count(*) as count").
		From(b.table).
		PlaceholderFormat(sq.Dollar)

	if filter != nil {
		qb = qb.Where(filter)
	}

	// Build SQL query
	query, args, err := qb.ToSql()
	if err != nil {
		return 0, err
	}

	// Prepare statement and defer it's closure
	stmt, err := b.session.PreparexContext(ctx, query)
	if err != nil {
		return 0, err
	}
	defer func(stmt *sqlx.Stmt) {
		//b.logger.SafeClose(stmt, "Unable to close statement")
	}(stmt)

	var count int64
	err = stmt.QueryRowxContext(ctx, args...).Scan(&count)
	if err != nil {
		return 0, err
	}

	// Return no error
	return count, nil
}

func (b *pgRepository) BeginTx(ctx context.Context, opts *sql.TxOptions) (*Tx, error) {
	beginx, err := b.session.BeginTxx(ctx, &sql.TxOptions{
		Isolation: sql.LevelDefault,
		ReadOnly:  false,
	})

	if err != nil {
		return nil, err
	}

	return newTx(beginx), nil
}

func (b *pgRepository) GetConn() *sqlx.DB {
	return b.session
}

func (b *pgRepository) WithTx(tx *Tx) PgTx {
	return newTxRepository(b.table, tx)
}

// -------------------------------------------------------------------------------------------

func (b *pgRepository) ExtractColumnPairs(data interface{}) ([]string, []interface{}) {
	// create type mapper
	valueMap := b.mapper.FieldMap(reflect.ValueOf(data))

	// Extract columns
	var columns = make([]string, len(valueMap))
	var values = make([]interface{}, len(valueMap))
	i := 0
	for column, value := range valueMap {
		columns[i] = column
		values[i] = value.Interface()
		i++
	}

	// Return all elements
	return columns, values
}
