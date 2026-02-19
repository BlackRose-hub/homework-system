package dao

import (
	"errors"                  //创建错误对象
	"homework-system/configs" //获取数据库连接配置
	"homework-system/models"  //数据模型
)

type HomeworkDao struct{} //创建作业数据访问对象的结构体

func NewHomeworkDao() *HomeworkDao {
	return &HomeworkDao{}
}

func (d *HomeworkDao) Create(homework *models.Homework) error { //数据库中创建作业记录
	return configs.DB.Create(homework).Error
}

func (d *HomeworkDao) FindByID(id uint) (*models.Homework, error) { //根据ID查询单个作业，并预加载创建者信息
	var homework models.Homework
	err := configs.DB.Preload("Creator").First(&homework, id).Error
	if err != nil {
		return nil, err
	}
	return &homework, nil
}

func (d *HomeworkDao) List(department string, page, pageSize int) ([]models.Homework, int64, error) { //分页查询作业列表，支持按部门筛选
	var homeworks []models.Homework
	var total int64 //记录作业数量

	db := configs.DB.Model(&models.Homework{}).Preload("Creator") //创建查询对象
	if department != "" {                                         //查询指定部门的作业
		db = db.Where("department = ?", department)
	}

	err := db.Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	err = db.Offset((page - 1) * pageSize).Limit(pageSize).Order("created_at desc").Find(&homeworks).Error //跳过前面的数据，只取一页的数据，按创建世家倒叙
	return homeworks, total, err
}

func (d *HomeworkDao) Update(homework *models.Homework) error { //更新作业信息，并处理并发
	// 版本号+1，并且只能更新版本号相同的记录
	result := configs.DB.Model(&models.Homework{}).
		Where("id = ? AND version = ?", homework.ID, homework.Version).
		Updates(map[string]interface{}{
			"title":       homework.Title,       //新标题
			"description": homework.Description, //新内容
			"deadline":    homework.Deadline,    //新截止时间
			"allow_late":  homework.AllowLate,   //新补交设置
			"version":     homework.Version + 1, //版本号+1
		})
	if result.Error != nil {
		return result.Error
	}
	if result.RowsAffected == 0 {
		return ErrVersionExpired
	} //没有更新记录，表示版本过期，被别人改过了
	return nil
}

func (d *HomeworkDao) Delete(id uint) error { //删除ID指定的作业
	return configs.DB.Delete(&models.Homework{}, id).Error
}

var ErrVersionExpired = errors.New("版本已过期，请刷新后重试") //定义版本过期的错误
