package db

import (
	"api/db/dbengine"
	"context"
	"fmt"
	"github.com/narvikd/errorskit"
)

func InsertFile(q *dbengine.Queries, fileName string, fileBytes []byte, chunkNum int64) error {
	_, errInsert := q.InsertFile(context.Background(), dbengine.InsertFileParams{
		Filename: fileName,
		Contents: fileBytes,
		ChunkNum: chunkNum,
	})
	if errInsert != nil {
		return fmt.Errorf("couldn't insert file to DB: '%s'. Err: %w", fileName, errInsert)
	}
	return nil
}

func GetFile(q *dbengine.Queries, fileName string, chunkNum int64) ([]byte, error) {
	bytes, errGet := q.GetFileByID(context.Background(), dbengine.GetFileByIDParams{
		Filename: fileName,
		ChunkNum: chunkNum,
	})
	if errGet != nil {
		return nil, errorskit.Wrap(errGet, "couldn't get file from db")
	}
	return bytes, nil
}
