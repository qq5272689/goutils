package service

import (
	"github.com/gin-gonic/gin"
	"github.com/qq5272689/goldden-go/pkg/models"
	"github.com/qq5272689/goldden-go/pkg/utils/crypto"
	"github.com/qq5272689/goldden-go/pkg/utils/http"
	"github.com/qq5272689/goldden-go/pkg/utils/logger"
	"github.com/qq5272689/goldden-go/pkg/utils/types"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

type UserService interface {
	GetUser(id int) (d models.User, err error)
	CheckPassword(id int, password string) (ok bool, err error)
	CreateUser(d *models.User) (err error)
	UpdateUser(d *models.User) (err error)
	DelUser(ids []int) (err error)
	SearchUser(filter string, pageNo, pageSize int) (td *types.TableData, err error)
}

type UserServiceDB struct {
	DB *gorm.DB
}

func GetUserServiceDB(db *gorm.DB) UserService {
	return &UserServiceDB{db}
}

func GetUserServiceDBWithContext(c *gin.Context) UserService {
	db, _ := c.Get("DB")
	return &UserServiceDB{db.(*gorm.DB)}
}

type GetUserError struct {
}

func (cue GetUserError) Error() string {
	return "GetUserError"
}

func (db *UserServiceDB) GetUser(id int) (d models.User, err error) {
	logger.Debug("GetUser 接受到任务：", zap.Int("id", id))
	tx := db.DB.Model(&d).
		Where(" id=?", id)
	err = tx.Last(&d).Error
	return
}

func (db *UserServiceDB) CheckPassword(id int, password string) (ok bool, err error) {
	logger.Debug("CheckPassword 接受到任务：", zap.Int("id", id))
	d := &models.User{}
	password = crypto.GetPassword(password)
	tx := db.DB.Model(d).
		Where(" id=? and password = ?", id, password)
	err = tx.Last(d).Error
	if err != nil {
		return false, err
	}
	if d.ID > 0 {
		return true, nil
	} else {
		return false, nil
	}
}

func (db *UserServiceDB) CreateUser(d *models.User) (err error) {
	logger.Debug("CreateUser 接受到任务：", zap.Reflect("args", *d))
	d.Password = crypto.GetPassword(d.Password)
	return db.DB.Create(d).Error
}

func (db *UserServiceDB) UpdateUser(d *models.User) (err error) {
	logger.Debug("UpdateUser 接受到任务：", zap.Reflect("args", *d))
	if d.Password != "" {
		d.Password = crypto.GetPassword(d.Password)
	}
	return db.DB.Model(&models.User{ID: d.ID}).Updates(d).Error
}

func (db *UserServiceDB) DelUser(ids []int) (err error) {
	logger.Debug("DelUser 接受到任务：", zap.Any("ids", ids))
	tx := db.DB.Begin()
	if err := tx.Where("id in ?", ids).Delete(&models.User{}).Error; err != nil {
		tx.Rollback()
		return err
	}
	tx.Commit()
	return nil
}

func (db *UserServiceDB) SearchUser(filter string, pageNo, pageSize int) (td *types.TableData, err error) {
	logger.Debug("SearchAlert接受到任务：", zap.String("filter", filter), zap.Int("pageno", pageNo), zap.Int("pagesize", pageSize))
	tx := db.DB.Model(&models.User{})
	if filter != "" {
		tx = tx.Where("content like ?", "%"+filter+"%")
	}
	var count int64
	if err = tx.Count(&count).Error; err != nil {
		return nil, err
	}
	tx.Limit(pageSize).Offset(pageSize * (pageNo - 1))
	ds := []models.User{}
	if err = tx.Find(&ds).Error; err != nil {
		return nil, err
	}

	return http.NewTableData(ds, pageNo, pageSize, int(count)), nil
}
