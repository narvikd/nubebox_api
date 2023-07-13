package route

import (
	"api/api/debugerr"
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
	const genericErr = "couldn't get filenames"

	f, err := c.Query.GetAllFileNames(context.Background())
	if err != nil {
		errorskit.LogWrap(err, genericErr)
		return jsonresponse.ServerError(fiberCtx, genericErr)
	}

	if len(f) <= 0 {
		return jsonresponse.ServerError(fiberCtx, "no files in DB")
	}

	return jsonresponse.OK(fiberCtx, "filenames retrieved from DB", f)
}

func (c *ApiCtx) replaceFile(fiberCtx *fiber.Ctx) error {
	file, errUpload := fiberCtx.FormFile("file")
	if errUpload != nil {
		debugMsg := debugerr.WrapMsg("process file upload", "uploadFile")
		errorskit.LogWrap(errUpload, debugMsg)
		return jsonresponse.BadRequest(fiberCtx, debugMsg)
	}

	errSave := fiberCtx.SaveFile(file, file.Filename)
	if errSave != nil {
		debugMsg := debugerr.WrapMsg("save file", "uploadFile")
		errorskit.LogWrap(errSave, debugMsg)
		return jsonresponse.ServerError(fiberCtx, debugMsg)
	}

	defer func() {
		if errDelete := filekit.DeleteFile(file.Filename); errDelete != nil {
			errorskit.LogWrap(errDelete, debugerr.WrapMsg("delete temporal file", "uploadVideo"))
		}
	}()

	_, errDeleteFile := c.Query.DeleteFile(context.Background(), file.Filename)
	if errDeleteFile != nil {
		debugMsg := debugerr.WrapMsg("delete old version of file", "uploadFile")
		errorskit.LogWrap(errDeleteFile, debugMsg)
		return jsonresponse.ServerError(fiberCtx, debugMsg)
	}

	if errFileToDB := operation.FileToDB(c.Query, file.Filename); errFileToDB != nil {
		debugMsg := debugerr.WrapMsg("save file to db", "uploadFile")
		errorskit.LogWrap(errFileToDB, debugMsg)
		return jsonresponse.ServerError(fiberCtx, debugMsg)
	}

	return jsonresponse.OK(fiberCtx, file.Filename, "")
}

func (c *ApiCtx) downloadFile(fiberCtx *fiber.Ctx) error {
	m := new(FileModel)
	if errParse := fiberparser.ParseAndValidate(fiberCtx, m); errParse != nil {
		return jsonresponse.BadRequest(fiberCtx, errParse.Error())
	}

	ext := filepath.Ext(m.FileName) // ext starts with a dot
	if ext == "" {
		return jsonresponse.BadRequest(fiberCtx, "filename was not valid")
	}

	destination := uuid.NewString() + ext

	errSave := operation.DBToFile(c.Query, m.FileName, destination)
	if errSave != nil {
		debugMsg := debugerr.WrapMsg("save file", "downloadFile")
		errorskit.LogWrap(errSave, debugMsg)
		return jsonresponse.ServerError(fiberCtx, debugMsg)
	}

	defer func() {
		if errDelete := filekit.DeleteFile(destination); errDelete != nil {
			errorskit.LogWrap(errDelete, debugerr.WrapMsg("delete temporal file", "downloadFile"))
		}
	}()

	return fiberCtx.Download(destination, m.FileName)
}

func (c *ApiCtx) deleteFile(fiberCtx *fiber.Ctx) error {
	m := new(FileModel)
	if errParse := fiberparser.ParseAndValidate(fiberCtx, m); errParse != nil {
		return jsonresponse.BadRequest(fiberCtx, errParse.Error())
	}

	exists, errCheckExists := c.Query.FileExists(context.Background(), m.FileName)
	if errCheckExists != nil {
		debugMsg := debugerr.WrapMsg("check if file exists", "deleteFile")
		errorskit.LogWrap(errCheckExists, debugMsg)
		return jsonresponse.ServerError(fiberCtx, debugMsg)
	}

	if !exists {
		return jsonresponse.BadRequest(fiberCtx, "file doesn't exist")
	}

	_, errDeleteFile := c.Query.DeleteFile(context.Background(), m.FileName)
	if errDeleteFile != nil {
		debugMsg := debugerr.WrapMsg("delete file", "deleteFile")
		errorskit.LogWrap(errDeleteFile, debugMsg)
		return jsonresponse.ServerError(fiberCtx, debugMsg)
	}

	return jsonresponse.OK(fiberCtx, m.FileName, "")
}
