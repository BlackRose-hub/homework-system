package dao

import (
	"homework-system/configs" //获取数据库连接
	"homework-system/models"  //使用User结构体
)

type UserDao struct{} //用户数据访问对象，负责相关的数据库操作

func NewUserDao() *UserDao { //创建UserDao实例
	return &UserDao{}
}
func (d *UserDao) Create(user *models.User) error { //创建用户
	return configs.DB.Create(user).Error //creat时GORM的插入方法
}
func (d *UserDao) FindByUsername(username string) (*models.User, error) { //根据用户名字查询用户信息
	var user models.User
	err := configs.DB.Where("username = ?", username).First(&user).Error //where设置查询条件，first获取第一条记录
	if err != nil {
		return nil, err //用户不存在或查询失败
	}
	return &user, nil
}
func (d *UserDao) FindByID(id uint) (*models.User, error) { //根据ID查询用户信息
	var user models.User
	err := configs.DB.First(&user, id).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
func (d *UserDao) Delete(id uint) error { //软删除用户
	return configs.DB.Delete(&models.User{}, id).Error
}
func (d *UserDao) CheckUsernameExist(username string) (bool, error) { //检查用户名是否已被注册
	var count int64
	err := configs.DB.Model(&models.User{}).Where("username = ?", username).Count(&count).Error //count统计数量
	return count > 0, err                                                                       //表示用户名已存在
}
