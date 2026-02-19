package models

import (
	"errors" //创建自定义错误
	"time"   //处理时间相关字段
)

type Homework struct { //定义作业表对应的GORM模型
	ID          uint       `gorm:"primaryKey" json:"id"`                                                                                 //作业唯一标识符，主键
	Title       string     `gorm:"type:varchar(200);nowt null" json:"title"`                                                             //作业标题，不能为空
	Description string     `gorm:"type:text;" json:"description"`                                                                        //作业详细描述，使用Text类型，可以存储长文本
	Department  Department `gorm:"type:enum('backend','frontend ','sre','product','design','android','ios');not null" json:"department"` //作业所属部门，不能为空，必须是七个部门之一
	CreatorID   uint       `gorm:"not null" json:"creator_id"`                                                                           //发布作业的老登的ID，不能为空
	Creator     User       `gorm:"foreignKey:CreatorID" json:"creator,omitempty"`                                                        //预加载的发布者信息，指定外键关系，为空则不序列化
	Deadline    time.Time  `json:"deadline"`                                                                                             //作业截至时间
	AllowLate   bool       `json:"allow_late"`                                                                                           //是否允许补交
	Version     int        `gorm:"default:0" json:"version"`                                                                             //版本，用以解决两个人同时修改同一份作业的情况，默认为0
	CreatedAt   time.Time  `json:"created_at"`                                                                                           //创建时间
	UpdatedAt   time.Time  `json:"updated_at"`                                                                                           //更新时间
}

func (h *Homework) GetDepartmentWithLabel() map[string]interface{} { //获取带中文标签的部门信息
	return map[string]interface{}{
		"department":       h.Department,                  //返回部门代码
		"department_label": DepartmentLabel[h.Department], //返回对应的中文名
	}

}
func (h *Homework) CanModify(user *User) bool {
	return user.Role == RoleAdmin && user.Department == h.Department
}
func (h *Homework) ValidateDeadline() error {
	if h.Deadline.Before(time.Now()) {
		return ErrDeadlinePassed
	}
	return nil
}

var (
	ErrDeadlinePassed = errors.New("截止时间不能早于当前时间")
)
