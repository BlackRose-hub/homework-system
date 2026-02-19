package models

import (
	"time" //时间字段处理

	"golang.org/x/crypto/bcrypt" //密码加密
	"gorm.io/gorm"               //ORM框架
)

type Role string       //角色类型
type Department string //部门类型

const ( //角色常量定义
	RoleStudent Role = "student" //小登
	RoleAdmin   Role = "admin"   //老登
)
const ( //部门常量定义
	DeptBackend Department = "backend" //后端
	DeptFronted Department = "fronted" //前端
	DepSRE      Department = "sre"     //运维
	DepProduct  Department = "product" //产品
	DeptDesign  Department = "design"  //设计
	DepAndroid  Department = "android" //Android
	DeptIOS     Department = "ios"     //IOS
)

var DepartmentLabel = map[Department]string{ //部门中文标签映射，用于API返回时显示中文部门名
	DeptBackend: "后端",
	DeptFronted: "前端",
	DepSRE:      "SRE",
	DepProduct:  "产品",
	DeptDesign:  "视觉设计",
	DepAndroid:  "Android",
	DeptIOS:     "IOS",
} //部门枚举值映射
// ·gorm：告诉数据库怎么存储该字段，json：告诉程序怎么返回给前端
type User struct { //User用户模型-对应数据库的user表
	ID         uint           `gorm:"primaryKey" json:"id"`                                                                               //主键                                                                            //主键
	Username   string         `gorm:"type:varchar(50);uniqueIndex" json:"username"`                                                       //用户名，唯一索引
	Password   string         `gorm:"type:varchar(255)" json:"-"`                                                                         //密码，不返回给前端
	Nickname   string         `gorm:"type:varchar(50)" json:"nickname"`                                                                   //用户昵称
	Role       Role           `gorm:"type:enum('student','admin');;default:'student'" json:"role"`                                        //身份只能是老登或者是小登，默认为小登
	Department Department     `gorm:"type:enum('backend','fronted','sre','product','design','android','ios');not null" json:"department"` //所属部门，只能是七个部门之一
	Email      string         `gorm:"type:varchar(100)" json:"email,omitempty"`                                                           //邮箱，omitempty=空时不返回
	CreatedAt  time.Time      `json:"created_at"`                                                                                         //创建时间
	UpdatedAt  time.Time      `json:"updated_at"`                                                                                         //更新时间
	DeletedAt  gorm.DeletedAt `gorm:"index" json:"-"`                                                                                     //软删除字段，不返回给前端
}

func (u *User) HashPassword() error { //对用户密码进行加密
	hashed, err := bcrypt.GenerateFromPassword([]byte(u.Password), bcrypt.DefaultCost)
	//bcrypt.Generate...对密码进行哈希加密
	//bcrypt.DefaultCost 加密强度
	if err != nil {
		return err //如果加密失败，返回错误
	}
	u.Password = string(hashed) //把加密后的密码存回去
	return nil
}
func (u *User) CheckPassword(password string) bool { //验证用户输入的密码是否匹配
	err := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(password)) //比较哈希值和明文密码
	return err == nil                                                          //如果错误为空，证明密码正确

}
func (u *User) GetDepartmentWithLabel() map[string]interface{} { //获取带中文标签的部门信息
	return map[string]interface{}{
		"department":       u.Department,                  //部门代码
		"department_label": DepartmentLabel[u.Department], //部门中文名
	}
}
