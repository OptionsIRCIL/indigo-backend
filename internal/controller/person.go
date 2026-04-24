package controller

import (
	"context"
	"crypto/sha512"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"myoptions.info/indigo/backend/internal/config"
	"myoptions.info/indigo/backend/internal/service"
	"myoptions.info/indigo/backend/internal/util"
	"myoptions.info/indigo/backend/model/entity"
)

func PersonAttachmentPost(database *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		token := util.FetchTokenFromContext(r)
		employee := util.TokenToEmployee(token, database, ctx)

		// Verify person exists
		parentId := r.PathValue("personId")
		parentCount, parentErr := gorm.G[entity.Person](database).Where("id = ?", parentId).Count(ctx, "id")
		if parentCount != 1 || parentErr != nil {
			util.ThrowHttpStatus(w, 404)
			return
		}

		// Given we're behind a front proxy, we won't worry about checking the attachment size.
		contentType := r.Header.Get("Content-Type")
		if !(contentType == "multipart/form-data" || strings.HasPrefix(contentType, "multipart/form-data;")) {
			util.ThrowHttpStatus(w, 415)
			return
		}

		formParseErr := r.ParseMultipartForm(int64(config.Config.Attachments.MaxFileSize))
		if formParseErr != nil {
			log.Println("Form parse error: ", formParseErr)
			util.ThrowHttpStatus(w, 422)
			return
		}

		file, fileOk := r.MultipartForm.File["attachment"]
		if !fileOk || file[0] == nil {
			log.Println("File error")
			util.ThrowHttpStatus(w, 422)
			return
		}

		fileContents, fileContentsErr := file[0].Open()
		if fileContentsErr != nil {
			log.Println("File contents error: ", fileContentsErr)
			util.ThrowHttpStatus(w, 422)
			return
		}

		fileBytes := make([]byte, file[0].Size)
		fileBytesCopied, fileBytesErr := fileContents.Read(fileBytes)

		if fileBytesCopied != int(file[0].Size) || fileBytesErr != nil {
			log.Println("File byte copy error: ", fileBytesCopied)
			util.ThrowHttpStatus(w, 422)
			return
		}

		fileId := uuid.New()
		mime, storeErr := service.StoreFile(fileId.String(), file[0].Filename, fileBytes)
		if storeErr != nil {
			log.Println("Encountered file storage error: ", storeErr)
			util.ThrowHttpStatus(w, 422)
			return
		}

		sha := sha512.New()
		sha.Write(fileBytes)
		sum := sha.Sum(nil)

		fileEntity := entity.PersonAttachment{
			Id:          fileId,
			EmployeeId:  employee.Id,
			PersonId:    uuid.MustParse(parentId),
			FileName:    file[0].Filename,
			ContentType: mime,
			Size:        uint(file[0].Size),
			Signature:   string(sum),
		}

		err := gorm.G[entity.PersonAttachment](database).Create(ctx, &fileEntity)
		if err != nil {
			// TODO: Unlink file?
			util.ThrowHttpUnhandled(w, err)
			return
		}

		util.ReturnSerialized(w, 200, fileEntity, []string{"get"})
	}
}

func PersonAttachmentGet(database *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()

		id := r.PathValue("id")
		count, err := gorm.G[entity.PersonAttachment](database).Where("id = ?", id).Count(ctx, "id")
		if count != 1 || err != nil {
			util.ThrowHttpStatus(w, 404)
			return
		}

		details, err := gorm.G[entity.PersonAttachment](database).Where("id = ?", id).First(ctx)
		if err != nil {
			util.ThrowHttpUnhandled(w, err)
			return
		}

		w.Header().Add("Content-Length", strconv.Itoa(int(details.Size)))
		w.Header().Add("Content-Type", details.ContentType)
		w.Header().Add("Content-Disposition", "attachment; filename=\""+url.QueryEscape(details.FileName)+"\"")

		if r.Method == "HEAD" {
			w.WriteHeader(200)
			return
		}

		contents, contentsErr := service.RetrieveFile(id)
		if contentsErr != nil {
			util.ThrowHttpUnhandled(w, contentsErr)
			return
		}

		w.Write(contents)
	}
}
