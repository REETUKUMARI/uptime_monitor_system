package controllers

import (
	"net/http"
	"net/url"
	"time"
	"website_status_checker/database"
	"website_status_checker/repository"

	"github.com/gin-gonic/gin"
	uuid "github.com/satori/go.uuid"
)

type p struct{}

type CreateUrlInput struct {
	URLLink          string        `json:"url" binding:"required"`
	CrawlTimeout     time.Duration `json:"crawl_timeout" binding:"required"`
	Frequency        int           `json:"frequency" binding:"required"`
	FailureThreshold int           `json:"failure_threshold" binding:"required"`
}

type UpdateUrlInput struct {
	CrawlTimeout     time.Duration `json:"crawl_timeout"`
	Frequency        int           `json:"frequency"`
	FailureThreshold int           `json:"failure_threshold"`
}

func GetUrls(c *gin.Context) {

	var url []database.Pingdom
	if err := repository.Repo.DatabaseGets(&url); err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		database.DB.Find(&url)
		c.JSON(http.StatusOK, url)
	}
}
func GetUrl(c *gin.Context) {
	id, err := uuid.FromString(c.Params.ByName("id"))
	res, err := repository.Repo.DatabaseGet(id)

	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, res)
	}
}

func CreateUrl(c *gin.Context) {
	var urll CreateUrlInput
	c.BindJSON(&urll)
	url, err := repository.Repo.DatabaseCreate(urll.URLLink, urll.CrawlTimeout, urll.Frequency, urll.FailureThreshold)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	} else {
		c.JSON(http.StatusOK, url)
	}
}
func Updateurl(c *gin.Context) {
	id, _ := uuid.FromString(c.Params.ByName("id"))
	var input database.Pingdom
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	url, err := repository.Repo.DatabaseUpdate(id, input.CrawlTimeout, input.Frequency, input.FailureThreshold)
	if err != nil {
		c.AbortWithStatus(http.StatusNotFound)
	}
	c.JSON(http.StatusOK, url)
}
func Deleteurl(c *gin.Context) {
	idd := c.Params.ByName("id")
	id, _ := uuid.FromString(idd)
	d := repository.Repo.DatabaseDelete(id)
	if d != nil {
		c.AbortWithStatus(http.StatusNotFound)
	}
	c.Writer.WriteHeader(http.StatusNoContent)
}
func IsUrl(str string) (bool, string) {
	u, err := url.Parse(str)
	return err == nil && u.Scheme != "" && u.Host != "", u.Scheme + "://" + u.Host + u.Path
}
