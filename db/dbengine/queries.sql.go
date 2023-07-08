// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.14.0
// source: queries.sql

package dbengine

import (
	"context"
	"database/sql"
)

const getFileByID = `-- name: GetFileByID :one
select contents
from testfiles
where filename = $1 and chunk_num = $2
limit 1
`

type GetFileByIDParams struct {
	Filename string `json:"filename"`
	ChunkNum int64  `json:"chunk_num"`
}

func (q *Queries) GetFileByID(ctx context.Context, arg GetFileByIDParams) ([]byte, error) {
	row := q.db.QueryRowContext(ctx, getFileByID, arg.Filename, arg.ChunkNum)
	var contents []byte
	err := row.Scan(&contents)
	return contents, err
}

const getFileChunkCount = `-- name: GetFileChunkCount :one
select count(chunk_num)
from testfiles
where filename = $1
limit 1
`

func (q *Queries) GetFileChunkCount(ctx context.Context, filename string) (int64, error) {
	row := q.db.QueryRowContext(ctx, getFileChunkCount, filename)
	var count int64
	err := row.Scan(&count)
	return count, err
}

const insertFile = `-- name: InsertFile :execresult
INSERT INTO testfiles (id, filename, contents, chunk_num)
VALUES (DEFAULT, $1, $2, $3)
`

type InsertFileParams struct {
	Filename string `json:"filename"`
	Contents []byte `json:"contents"`
	ChunkNum int64  `json:"chunk_num"`
}

func (q *Queries) InsertFile(ctx context.Context, arg InsertFileParams) (sql.Result, error) {
	return q.db.ExecContext(ctx, insertFile, arg.Filename, arg.Contents, arg.ChunkNum)
}