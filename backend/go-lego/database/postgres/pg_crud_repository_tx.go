package postgres

import (
	"context"
	"fmt"
	"reflect"

	"github.com/thiagoretondar/golang-blog-example/backend/go-lego/database"

	sq "github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	"github.com/jmoiron/sqlx/reflectx"
)

type PgTx interface {
	database.CRUDRepository
	GetTxConn() *sqlx.Tx
}

type PgTxRepository struct {
	table  string
	tx     *sqlx.Tx
	mapper *reflectx.Mapper
}

func newTxRepository(tableName string, session *Tx) PgTx {
	return &PgTxRepository{
		table:  tableName,
		tx:     session.getTx(),
		mapper: reflectx.NewMapper("db"),
	}
}

// Insert inserts a single record
func (b *PgTxRepository) Insert(ctx context.Context, data interface{}, lastInsertedID interface{}) error {
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

	// The statements prepared for a transaction by calling the transaction's Prepare or Stmt methods
	// are closed by the call to Commit or Rollback.
	stmt, err := b.tx.PreparexContext(ctx, query)
	if err != nil {
		return err
	}

	// Do the insert query
	if lastInsertedID != nil {
		// We do QueryRowxContext because Postgres doesn't work with lastInsertedID
		err = stmt.QueryRowxContext(ctx, args...).Scan(lastInsertedID)
	} else {
		// Here we don't need the lastInsertedID
		_, err = stmt.ExecContext(ctx, args...)
	}

	return err
}

// FindAll returns all records from database
func (b *PgTxRepository) FindAll(ctx context.Context, output interface{}) error {
	//TODO implement pagination

	// Prepare query
	queryBuilder := sq.Select("*").From(b.table).PlaceholderFormat(sq.Dollar)

	// Build SQL Query
	query, _, err := queryBuilder.ToSql()
	if err != nil {
		return fmt.Errorf("unable to build query: %w", err)
	}

	// The statements prepared for a transaction by calling the transaction's Prepare or Stmt methods
	// are closed by the call to Commit or Rollback.
	stmt, err := b.tx.PreparexContext(ctx, query)
	if err != nil {
		return err
	}

	err = stmt.SelectContext(ctx, output)
	if err != nil {
		return err
	}
	return nil
}

// FindOne returns only one record given the filter
func (b *PgTxRepository) FindOne(ctx context.Context, filter interface{}, output interface{}) error {
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

	// The statements prepared for a transaction by calling the transaction's Prepare or Stmt methods
	// are closed by the call to Commit or Rollback.
	stmt, err := b.tx.PreparexContext(ctx, query)
	if err != nil {
		return err
	}

	// Do the find
	err = stmt.QueryRowxContext(ctx, args...).StructScan(output)
	if err != nil {
		return err
	}

	// No errors
	return nil
}

// Find returns all records from database that match the filter
func (b *PgTxRepository) Find(ctx context.Context, filter interface{}, output interface{}) error {
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

	// The statements prepared for a transaction by calling the transaction's Prepare or Stmt methods
	// are closed by the call to Commit or Rollback.
	stmt, err := b.tx.PreparexContext(ctx, query)
	if err != nil {
		return err
	}

	// Do the find
	err = stmt.SelectContext(ctx, output, args...)
	if err != nil {
		return err
	}

	// No errors
	return nil
}

// Update updates records matching the given filter
func (b *PgTxRepository) Update(ctx context.Context, set map[string]interface{}, filter interface{}) (int64, error) {
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

	// The statements prepared for a transaction by calling the transaction's Prepare or Stmt methods
	// are closed by the call to Commit or Rollback.
	stmt, err := b.tx.PreparexContext(ctx, query)
	if err != nil {
		return 0, err
	}

	// Execute query
	result, err := stmt.ExecContext(ctx, args...)
	if err != nil {
		return 0, err
	}

	// Return result and possible error
	return result.RowsAffected()
}

// Remove updates the records that match the given filter to deleted status. The record isn't really deleted from database
func (b *PgTxRepository) Remove(ctx context.Context, filter interface{}, physicalDeletion bool) (int64, error) {
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

	// The statements prepared for a transaction by calling the transaction's Prepare or Stmt methods
	// are closed by the call to Commit or Rollback.
	stmt, err := b.tx.PreparexContext(ctx, query)
	if err != nil {
		return 0, err
	}

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
func (b *PgTxRepository) Count(ctx context.Context, filter interface{}) (int64, error) {
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

	// The statements prepared for a transaction by calling the transaction's Prepare or Stmt methods
	// are closed by the call to Commit or Rollback.
	stmt, err := b.tx.PreparexContext(ctx, query)
	if err != nil {
		return 0, err
	}

	var count int64
	err = stmt.QueryRowxContext(ctx, args...).Scan(&count)
	if err != nil {
		return 0, err
	}

	// Return no error
	return count, nil
}

func (b *PgTxRepository) ExtractColumnPairs(data interface{}) ([]string, []interface{}) {
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

func (b *PgTxRepository) GetTxConn() *sqlx.Tx {
	return b.tx
}
