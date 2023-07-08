package operation

import (
	"api/db"
	"api/db/dbengine"
	"context"
	"errors"
	"fmt"
	"github.com/narvikd/errorskit"
	"github.com/narvikd/filekit"
	"io"
	"log"
	"os"
	"path/filepath"
)

func DBToFile(query *dbengine.Queries, srcFileName string, dstFileName string) error {
	var (
		i         int64
		fileNames []string
	)

	count, errChunkCount := query.GetFileChunkCount(context.Background(), srcFileName)
	if errChunkCount != nil {
		return errorskit.Wrap(errChunkCount, "couldn't get file chunk count")
	}

	if count <= 0 {
		return errors.New("file doesn't exist on DB")
	}

	for i = 0; i < count; i++ {
		fileBytes, errGetFile := db.GetFile(query, srcFileName, i)
		if errGetFile != nil {
			return errGetFile
		}

		chunkFilePath := filepath.Join(
			filepath.Dir(dstFileName),
			fmt.Sprintf("%s_chunk_%d",
				filepath.Base(dstFileName), i,
			),
		)
		fileNames = append(fileNames, chunkFilePath)

		errWrite := filekit.WriteToFile(chunkFilePath, fileBytes)
		if errWrite != nil {
			return errWrite
		}
	}

	defer func(fileNames []string) {
		if err := deleteChunkFiles(fileNames); err != nil {
			log.Println(err)
		}
	}(fileNames)

	if errJoinChunks := joinChunks(fileNames, dstFileName); errJoinChunks != nil {
		return errJoinChunks
	}

	return nil
}

// joinChunks takes a slice of filenames and a destination file path as input.
// It then reads each chunk file in order and writes its contents to the destination file.
func joinChunks(fileNames []string, destFilePath string) error {
	destFile, errCreate := os.Create(filepath.Clean(destFilePath))
	if errCreate != nil {
		return errCreate
	}
	defer destFile.Close()

	for _, chunkFilePath := range fileNames {
		chunkFile, errOpen := os.Open(filepath.Clean(chunkFilePath))
		if errOpen != nil {
			return errOpen
		}

		_, errCopy := io.Copy(destFile, chunkFile)
		if errCopy != nil {
			// TODO: Refactor
			_ = chunkFile.Close()
			return errCopy
		}

		_ = chunkFile.Close()
	}

	return nil
}
