package tracks

import (
	"net/http"
	"time"

	"github.com/antongoncharik/crypto-knight-api/internal/entity/track"
	"github.com/antongoncharik/crypto-knight-api/internal/service"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

type Tracks struct {
	svc *service.Service
}

func New(svc *service.Service) *Tracks {
	return &Tracks{svc}
}

func (t *Tracks) GetAll(ctx *gin.Context) {
	var queryParams track.QueryParams

	err := ctx.ShouldBindQuery(&queryParams)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	now := time.Now()

	if queryParams.From == "" {
		startOfDay := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location())
		queryParams.From = startOfDay.Format("2006-01-02 15:04:05")
	}
	if queryParams.To == "" {
		endOfDay := time.Date(now.Year(), now.Month(), now.Day(), 23, 59, 59, 0, now.Location())
		queryParams.To = endOfDay.Format("2006-01-02 15:04:05")
	}

	_, err = time.Parse("2006-01-02 15:04:05", queryParams.From)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Expected 'YYYY-MM-DD HH:MI:SS'"})
		return
	}
	_, err = time.Parse("2006-01-02 15:04:05", queryParams.To)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Expected 'YYYY-MM-DD HH:MI:SS'"})
		return
	}

	res, err := t.svc.Tracks.GetAll(queryParams)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, res)
}

func (t *Tracks) Create(ctx *gin.Context) {
	var trackData track.Track

	err := ctx.ShouldBindJSON(&trackData)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	validate := validator.New()
	err = validate.Struct(trackData)
	if err != nil {
		errors := make(map[string]string)
		for _, e := range err.(validator.ValidationErrors) {
			errors[e.Field()] = e.Tag()
		}
		ctx.JSON(http.StatusBadRequest, gin.H{"error": errors})
		return
	}

	err = t.svc.Tracks.Create(trackData)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.Status(http.StatusCreated)
}
