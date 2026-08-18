package handler

import (
	"context"
	"net/http"
	"os"

	"github.com/CogitoNTNU/cogi-go/internal/service"
	"github.com/gin-gonic/gin"
)

type S3Tester struct {
	service *service.S3
	ctx     *context.Context
}

func NewS3Tester(s3Service *service.S3, ctx *context.Context) *S3Tester {
	return &S3Tester{service: s3Service, ctx: ctx}
}

func (s *S3Tester) TestUpload(gCtx *gin.Context) {
	path := "/home/adrianjlund/NTNU/Cogito/download (6).jpg"

	data, err := os.ReadFile(path)
	if err != nil {
		gCtx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":  "failed to read file",
			"detail": err.Error(),
		})
		return
	}

	err = s.service.UploadFile(*s.ctx, "test-image.jpg", data)
	if err != nil {
		gCtx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	url := s.service.GetPublicURL("test-image.jpg")

	gCtx.JSON(http.StatusOK, gin.H{
		"message": "File uploaded successfully",
		"url":     url,
	})
}

func (s *S3Tester) TestUploadLarge(gCtx *gin.Context) {
	path := "/home/adrianjlund/NTNU/Cogito/download (6).jpg"

	data, err := os.ReadFile(path)
	if err != nil {
		gCtx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error":  "failed to read file",
			"detail": err.Error(),
		})
		return
	}

	err = s.service.UploadLargeFile(*s.ctx, "test-large-image.jpg", data)
	if err != nil {
		gCtx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	url := s.service.GetPublicURL("test-large-image.jpg")

	gCtx.JSON(http.StatusOK, gin.H{
		"message": "Large file uploaded successfully",
		"url":     url,
	})
}

func (s *S3Tester) TestDelete(gCtx *gin.Context) {
	err := s.service.DeleteFile(*s.ctx, "test-image.jpg")
	if err != nil {
		gCtx.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	gCtx.JSON(http.StatusOK, gin.H{
		"message": "File deleted successfully",
	})
}