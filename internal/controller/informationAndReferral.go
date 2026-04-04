package controller

import (
	"context"
	"crypto/sha512"
	"log"
	"net/http"
	"strings"

	"github.com/google/uuid"
	"gorm.io/gorm"
	"myoptions.info/indigo/backend/internal/config"
	"myoptions.info/indigo/backend/internal/service"
	"myoptions.info/indigo/backend/internal/util"
	"myoptions.info/indigo/backend/model/entity"
)

func InformationAndReferralEffortPost(database *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		token := util.FetchTokenFromContext(r)
		employee := util.TokenToEmployee(token, database, ctx)

		// TODO: Users currently have the ability to create effort records for other employees if
		//       the employeeId key is specified. This should be restricted to only be allowed for administrators.
		deserializationErr, deserialized := util.Deserialize[entity.InformationAndReferralEffort](r.Body, []string{"post"})
		if deserializationErr != nil {
			// TODO: Are these messages safe to relay back to client?
			util.ThrowHttpError(w, 422, "Could not deserialize POST body: "+deserializationErr.Error())
			return
		}

		deserialized.InformationAndReferralId, _ = uuid.Parse(r.PathValue("informationAndReferralId"))

		// Set employee ID if not already provided
		if deserialized.EmployeeId.String() == "00000000-0000-0000-0000-000000000000" {
			deserialized.EmployeeId = employee.Id
		}

		createErr := gorm.G[entity.InformationAndReferralEffort](database).Create(ctx, &deserialized)
		if createErr != nil {
			util.ThrowHttpUnhandled(w, createErr)
			return
		}

		util.ReturnSerialized(w, 201, deserialized, []string{"get"})
	}
}

func InformationAndReferralAttachmentPost(database *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		token := util.FetchTokenFromContext(r)
		employee := util.TokenToEmployee(token, database, ctx)

		// Verify I&R exists
		parentId := r.PathValue("informationAndReferralId")
		parentCount, parentErr := gorm.G[entity.InformationAndReferral](database).Where("id = ?", parentId).Count(ctx, "id")
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
			util.ThrowHttpUnhandled(w, storeErr)
			return
		}

		sha := sha512.New()
		sha.Write(fileBytes)
		sum := sha.Sum(nil)

		fileEntity := entity.InformationAndReferralAttachment{
			Id:                       fileId,
			EmployeeId:               employee.Id,
			InformationAndReferralId: uuid.MustParse(parentId),
			FileName:                 file[0].Filename,
			ContentType:              mime,
			Size:                     uint(file[0].Size),
			Signature:                string(sum),
		}

		err := gorm.G[entity.InformationAndReferralAttachment](database).Create(ctx, &fileEntity)
		if err != nil {
			// TODO: Unlink file?
			util.ThrowHttpUnhandled(w, err)
			return
		}

		util.ReturnSerialized(w, 200, fileEntity, []string{"get"})
	}
}
