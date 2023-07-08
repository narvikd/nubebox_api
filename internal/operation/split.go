package operation

import (
	"api/db"
	"api/db/dbengine"
	"errors"
	"fmt"
	"github.com/narvikd/errorskit"
	"github.com/narvikd/filekit"
	"io"
	"log"
	"os"
	"path/filepath"
)

func FileToDB(query *dbengine.Queries, fileName string) error {
	const chunkSize = 1 << 20 // 1MB
	fileList, errSplit := splitFileIntoChunks(fileName, chunkSize)
	if errSplit != nil {
		return errSplit
	}

	defer func(fileList []string) {
		if err := deleteChunkFiles(fileList); err != nil {
			log.Println(err)
		}
	}(fileList)

	for i, file := range fileList {
		fileBytes, errRead := filekit.ReadFile(file)
		if errRead != nil {
			return errRead
		}

		errInsert := db.InsertFile(query, fileName, fileBytes, int64(i))
		if errInsert != nil {
			return errInsert
		}
		log.Println("processed:", file)
	}

	return nil
}

// splitFileIntoChunks is a function that splits the provided file into chunks of a given size.
// The function returns a slice of strings containing the file names of each chunk file.
func splitFileIntoChunks(fileName string, chunkSize int) ([]string, error) {
	fileDescriptor, errOpen := os.Open(filepath.Clean(fileName))
	if errOpen != nil {
		return nil, errorskit.Wrap(errOpen, "couldn't open file at splitter")
	}
	defer func(fileDescriptor *os.File) {
		if errCloser := fileDescriptor.Close(); errCloser != nil {
			log.Printf("warning: OS couldn't close file: '%s' correctly. Err: %v\n", fileName, errCloser)
		}
	}(fileDescriptor)

	var fileNames []string
	chunkNumCounter := -1

	// Infinite loop to read the file in chunks until it's completely read
	for {
		chunkNumCounter++
		chunkBuffer := make([]byte, chunkSize)

		bytesRead, errRead := fileDescriptor.Read(chunkBuffer)
		if errRead != nil && errRead != io.EOF {
			return nil, errorskit.Wrap(errRead, "couldn't read chunk at splitter")
		}

		if bytesRead == 0 {
			break
		}

		// Create a unique file name for each chunk
		chunkFilePath := filepath.Join(
			filepath.Dir(fileName),
			fmt.Sprintf("%s_chunk_%d",
				filepath.Base(fileName), chunkNumCounter,
			),
		)

		chunkFile, errCreate := os.Create(filepath.Clean(chunkFilePath))
		if errCreate != nil {
			return nil, errorskit.Wrap(errCreate, "couldn't create chunked file at splitter")
		}

		_, errWrite := chunkFile.Write(chunkBuffer[:bytesRead]) // Only write the bytes that were read
		if errWrite != nil {
			// It is safe to ignore this error, as the size of the file will be enough for the kernel to handle
			// even in very slow drives (1mb/s)
			_ = chunkFile.Close()
			return nil, errorskit.Wrap(errWrite, "couldn't save chunked file to disk at splitter")
		}
		// Same as above, this is not deferred as we're in a loop, and it could very well be executed in parallel
		_ = chunkFile.Close()

		fileNames = append(fileNames, chunkFilePath)
	}

	if len(fileNames) == 0 {
		return nil, errors.New("splitter didn't return a slice of filenames")
	}

	return fileNames, nil
}
