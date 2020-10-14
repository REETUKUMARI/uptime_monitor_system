package monitor

import (
	"database/sql"
	"website_status_checker/database"

	uuid "github.com/satori/go.uuid"
)

type MonitorReposController interface {
	DatabaseSave(id uuid.UUID)
	DatabaseGetURL(urllink string) (database.Pingdom, error)
	DatabaseSaveFailureCount(urllink string)
	IncreaseFailureCount(id uuid.UUID)
	FailureCountToZero(id uuid.UUID)
	UpdateStatus(id uuid.UUID, st string)
	GetUrlData(urllink string) database.Pingdom
	GetRows() (*sql.Rows, error)
}
type MonitorRepoService struct{}

var (
	MonitorRepo MonitorReposController = MonitorRepoService{}
)

func (rp MonitorRepoService) DatabaseSave(id uuid.UUID) {
	var p database.Pingdom
	database.DB.First(&p, "id  = ?", id)
	database.DB.Model(&p).Update("FailureCount", p.FailureCount+1)
	database.DB.Model(&p).Update("Status", "inactive")
}
func (rp MonitorRepoService) DatabaseGetURL(urllink string) (database.Pingdom, error) {
	var url database.Pingdom
	result := database.DB.First(&url, "url_link  = ?", urllink)
	if result.Error != nil {
		return database.Pingdom{}, result.Error
	}
	return url, nil
}
func (rp MonitorRepoService) DatabaseSaveFailureCount(urllink string) {
	var p database.Pingdom
	database.DB.First(&p, "url_link  = ?", urllink)
	database.DB.Model(&p).Update("FailureCount", 0)
}
func (rp MonitorRepoService) IncreaseFailureCount(id uuid.UUID) {
	var url database.Pingdom
	database.DB.Where("id = ?", id).First(&url)

	database.DB.Model(&url).Update("FailureCount", url.FailureCount+1)
}
func (rp MonitorRepoService) FailureCountToZero(id uuid.UUID) {
	var url database.Pingdom
	database.DB.Where("id = ?", id).First(&url)

	database.DB.Model(&url).Update("FailureCount", 0)
}

func (rp MonitorRepoService) UpdateStatus(id uuid.UUID, st string) {
	var url database.Pingdom
	database.DB.Where("id = ?", id).First(&url)

	database.DB.Model(&url).Update("Status", st)
}
func (rp MonitorRepoService) GetUrlData(urllink string) database.Pingdom {
	var p database.Pingdom
	database.DB.First(&p, "url_link  = ?", urllink)
	return p
}
func (rp MonitorRepoService) GetRows() (*sql.Rows, error) {
	return database.DB.Raw("select * from pingdoms WHERE status = ?", "inactive").Rows()
}
