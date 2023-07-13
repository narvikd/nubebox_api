package operation

import "github.com/narvikd/filekit"

func deleteChunkFiles(fileNames []string) error {
	for _, v := range fileNames {
		if err := filekit.DeleteFile(v); err != nil {
			return err
		}
	}
	return nil
}
