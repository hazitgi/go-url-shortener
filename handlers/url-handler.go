package handlers

import (
	"fmt"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/hazi-tgi/go-url-shortner/controllers"
	"github.com/hazi-tgi/go-url-shortner/utils"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type UrlHandler struct {
	apiPrefix     string
	urlController controllers.UrlController
}

func NewUrlHandler(urlManager controllers.UrlController) *UrlHandler {
	return &UrlHandler{
		"api/v1",
		urlManager,
	}
}

func (h *UrlHandler) RegisterRoutes(r *gin.Engine) {
	api := r.Group(h.apiPrefix)
	{
		api.GET("/urls", h.findAll)
		api.GET("/make-short", h.makeShort)
	}
	r.GET("/:id", h.redirectToOriginalUrl)
}

func (h *UrlHandler) findAll(ctx *gin.Context) {
	data, err := h.urlController.FindAll()
	if err != nil {
		ctx.JSON(500, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(200, gin.H{"data": data})
}

func (h *UrlHandler) makeShort(ctx *gin.Context) {
	url := ctx.Query("url")
	fmt.Println("url", url)
	shortDomain := os.Getenv("SHORT_URI_DOMAIN")
	if url == "" {
		fmt.Println("URL is empty")
		response := utils.HTTPError{
			Success: false,
			Message: "URL is empty",
			Error:   "URL is empty",
			Status:  500,
		}
		response.InValidResponse(ctx)
		return
	}
	shortenedUrl, err := h.urlController.MakeShort(url)
	if err != nil {
		response := utils.HTTPError{
			Success: false,
			Message: "Error in making short url",
			Error:   err.Error(),
			Status:  500,
		}
		response.InValidResponse(ctx)
		return
	}

	response := utils.HTTPSuccess{
		Success: true,
		Message: "Success",
		Data: map[string]interface{}{
			"short_url": shortDomain + "/" + shortenedUrl.ID.Hex(),
		},
		Status: 200,
	}
	response.SuccessResponse(ctx)
}

func (h *UrlHandler) redirectToOriginalUrl(ctx *gin.Context) {
	id := ctx.Param("id")
	fmt.Println(">>>>>>> id", id)
	objectId, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		response := utils.HTTPError{
			Success: false,
			Message: "Invalid URL",
			Error:   err.Error(),
			Status:  500,
		}
		response.InValidResponse(ctx)
		return
	}
	url, err := h.urlController.FindById(objectId)
	if err != nil {
		response := utils.HTTPError{
			Success: false,
			Message: "Error in finding url",
			Error:   err.Error(),
			Status:  500,
		}
		response.InValidResponse(ctx)
		return
	}

	ctx.Redirect(301, url.Url)
}
