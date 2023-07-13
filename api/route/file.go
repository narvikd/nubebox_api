package route

import (
	"api/api/jsonresponse"
	"api/internal/operation"
	"context"
	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/narvikd/errorskit"
	"github.com/narvikd/fiberparser"
	"github.com/narvikd/filekit"
	"path/filepath"
)

type FileModel struct {
	FileName string `json:"filename" validate:"required"`
}

func (c *ApiCtx) listFiles(fiberCtx *fiber.Ctx) error {
	f, err := c.Query.GetAllFileNames(context.Background())
	if err != nil {
		const clientErr = "couldn't get filenames"
		errorskit.LogWrap(err, clientErr)
		return jsonresponse.ServerError(fiberCtx, clientErr)
	}

	if len(f) <= 0 {
		return jsonresponse.ServerError(fiberCtx, "no files in DB")
	}

	return jsonresponse.OK(fiberCtx, "filenames retrieved from DB", f)
}

func (c *ApiCtx) replaceFile(fiberCtx *fiber.Ctx) error {
	file, errUpload := fiberCtx.FormFile("file")
	if errUpload != nil {
		wrpErr := errorskit.Wrap(errUpload, "couldn't process formFile at uploadFile")
		return jsonresponse.BadRequest(fiberCtx, wrpErr.Error())
	}
	errSave := fiberCtx.SaveFile(file, file.Filename)
	if errSave != nil {
		wrpErr := errorskit.Wrap(errSave, "couldn't save file at uploadFile")
		return jsonresponse.ServerError(fiberCtx, wrpErr.Error())
	}
	defer func() {
		if errDelete := filekit.DeleteFile(file.Filename); errDelete != nil {
			errorskit.LogWrap(errDelete, "couldn't delete temporal file at uploadVideo")
		}
	}()

	_, errDeleteFile := c.Query.DeleteFile(context.Background(), file.Filename)
	if errDeleteFile != nil {
		return jsonresponse.ServerError(fiberCtx, errDeleteFile.Error())
	}

	if errFileToDB := operation.FileToDB(c.Query, file.Filename); errFileToDB != nil {
		wrpErr := errorskit.Wrap(errFileToDB, "couldn't save file to db at uploadVideo")
		return jsonresponse.ServerError(fiberCtx, wrpErr.Error())
	}
	return jsonresponse.OK(fiberCtx, file.Filename, "")
}

func (c *ApiCtx) downloadFile(fiberCtx *fiber.Ctx) error {
	m := new(FileModel)
	errParse := fiberparser.ParseAndValidate(fiberCtx, m)
	if errParse != nil {
		return jsonresponse.BadRequest(fiberCtx, errParse.Error())
	}
	ext := filepath.Ext(m.FileName)
	if ext == "" {
		return jsonresponse.BadRequest(fiberCtx, "filename was not valid")
	}

	destination := uuid.NewString() + ext

	errSaveFile := operation.DBToFile(c.Query, m.FileName, destination)
	if errSaveFile != nil {
		return jsonresponse.ServerError(fiberCtx, errSaveFile.Error())
	}
	defer func() {
		if errDelete := filekit.DeleteFile(destination); errDelete != nil {
			errorskit.LogWrap(errDelete, "couldn't delete temporal file at downloadFile")
		}
	}()

	return fiberCtx.Download(destination, m.FileName)
}

func (c *ApiCtx) deleteFile(fiberCtx *fiber.Ctx) error {
	m := new(FileModel)
	errParse := fiberparser.ParseAndValidate(fiberCtx, m)
	if errParse != nil {
		return jsonresponse.BadRequest(fiberCtx, errParse.Error())
	}

	exists, errCheckExists := c.Query.FileExists(context.Background(), m.FileName)
	if errCheckExists != nil {
		return jsonresponse.ServerError(fiberCtx, errCheckExists.Error())
	}

	if !exists {
		return jsonresponse.BadRequest(fiberCtx, "file doesn't exist")
	}

	_, errDeleteFile := c.Query.DeleteFile(context.Background(), m.FileName)
	if errDeleteFile != nil {
		return jsonresponse.ServerError(fiberCtx, errDeleteFile.Error())
	}

	return jsonresponse.OK(fiberCtx, m.FileName, "")
}
