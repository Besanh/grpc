package database

import (
	"crypto/tls"
	"database/sql"
	"time"

	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
)

type (
	IPostgreSql interface {
		Connect() (*bun.DB, error)
		GetDB() *bun.DB
	}
	PostgreSql struct {
		Driver       string
		Host         string
		Port         int
		Username     string
		Password     string
		Database     string
		Timeout      int
		DialTimeout  int
		ReadTimeout  int
		WriteTimeout int
		MaxIdleConns int
		MaxOpenConns int
	}
	PostgreSqlCon struct {
		PostgreSql
		DB *bun.DB
	}
)

func NewPostgreSqlCon(postgreSql PostgreSql) *PostgreSqlCon {
	pgCon := PostgreSqlCon{
		PostgreSql: postgreSql,
	}

	db, err := pgCon.Connect()
	if err != nil {
		return nil
	}
	pgCon.DB = db
	if err := pgCon.DB.Ping(); err != nil {
		return nil
	}
	return &pgCon
}

func (p *PostgreSqlCon) Connect() (*bun.DB, error) {
	pgconn := pgdriver.NewConnector(
		pgdriver.WithNetwork("tcp"),
		pgdriver.WithAddr(p.Host+":"+string(rune(p.Port))),
		pgdriver.WithTLSConfig(&tls.Config{InsecureSkipVerify: true}),
		pgdriver.WithUser(p.Username),
		pgdriver.WithPassword(p.Password),
		pgdriver.WithDatabase(p.Database),
		pgdriver.WithTimeout(time.Duration(p.Timeout)*time.Second),
		pgdriver.WithDialTimeout(time.Duration(p.DialTimeout)*time.Second),
		pgdriver.WithReadTimeout(time.Duration(p.ReadTimeout)*time.Second),
		pgdriver.WithWriteTimeout(time.Duration(p.WriteTimeout)*time.Second),
	)
	sqldb := sql.OpenDB(pgconn)
	sqldb.SetMaxIdleConns(p.MaxIdleConns)
	sqldb.SetMaxOpenConns(p.MaxOpenConns)
	bun.NewDB(sqldb, pgdialect.New())

	return bun.NewDB(sqldb, pgdialect.New()), nil
}

func (p *PostgreSqlCon) GetDB() *bun.DB {
	return p.DB
}
