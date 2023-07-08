package operation

import "github.com/narvikd/filekit"

func deleteChunkFiles(fileNames []string) error {
	for _, v := range fileNames {
		err := filekit.DeleteFile(v)
		if err != nil {
			return err
		}
	}
	return nil
}
