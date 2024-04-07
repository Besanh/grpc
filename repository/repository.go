package repository

import (
	"anhle-grpc/lib/database"
	"anhle-grpc/lib/elasticsearch"
	"context"
)

var DBCON database.IPostgreSql
var ESCON elasticsearch.IElasticsearch

func CreateTable(ctx context.Context, db database.IPostgreSql, entity any) (err error) {
	_, err = db.GetDB().NewCreateTable().Model(entity).
		IfNotExists().
		Exec(ctx)
	return
}

func InitRepos() {

}

func InitTables(ctx context.Context, db database.PostgreSql) {

}

func InitColumns() {

}

func InitIndex() {

}
