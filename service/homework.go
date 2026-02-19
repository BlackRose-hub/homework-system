package service

import (
	"fmt"
	"homework-system/configs"
	"homework-system/dao"
	"homework-system/models"
	"homework-system/pkg/errcode"
	"log"
	"time"
)

type HomeworkService struct {
	homeworkDao *dao.HomeworkDao
	userDao     *dao.UserDao
}

func NewHomeworkService() *HomeworkService {
	return &HomeworkService{
		homeworkDao: dao.NewHomeworkDao(),
		userDao:     dao.NewUserDao(),
	}
}

// 创建作业请求体
type CreateHomeworkRequest struct {
	Title       string    `json:"title" binding:"required,min=1,max=200"`
	Description string    `json:"description" binding:"required,min=1"`
	Department  string    `json:"department" binding:"required,oneof=backend frontend sre product design android ios"`
	Deadline    time.Time `json:"deadline" binding:"required,gt"`
	AllowLate   bool      `json:"allow_late"`
}

// 更新作业请求体
type UpdateHomeworkRequest struct {
	ID          uint      `json:"id" binding:"required"`
	Title       string    `json:"title" binding:"required,min=1,max=200"`
	Description string    `json:"description" binding:"required,min=1"`
	Deadline    time.Time `json:"deadline" binding:"required,gt"`
	AllowLate   bool      `json:"allow_late"`
	Version     int       `json:"version" binding:"required,min=0"`
}

// 作业列表响应
type HomeworkListResponse struct {
	Total int64                  `json:"total"`
	Items []HomeworkItemResponse `json:"items"`
	Page  int                    `json:"page"`
	Size  int                    `json:"size"`
}

type HomeworkItemResponse struct {
	ID              uint      `json:"id"`
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	Department      string    `json:"department"`
	DepartmentLabel string    `json:"department_label"`
	Creator         Creator   `json:"creator"`
	Deadline        time.Time `json:"deadline"`
	AllowLate       bool      `json:"allow_late"`
	Version         int       `json:"version"`
	CreatedAt       time.Time `json:"created_at"`
}

type Creator struct {
	ID              uint   `json:"id"`
	Nickname        string `json:"nickname"`
	Department      string `json:"department"`
	DepartmentLabel string `json:"department_label"`
}

// 发布作业（带极简并发示例）
func (s *HomeworkService) Create(creatorID uint, req *CreateHomeworkRequest) (*models.Homework, int, error) {
	// 验证发布者是否为管理员
	creator, err := s.userDao.FindByID(creatorID)
	if err != nil {
		return nil, errcode.UserNotFound, err
	}
	if creator.Role != models.RoleAdmin {
		return nil, errcode.PermissionDenied, nil
	}

	// 验证截止时间
	if req.Deadline.Before(time.Now()) {
		return nil, errcode.InvalidParam, models.ErrDeadlinePassed
	}

	// 创建作业
	homework := &models.Homework{
		Title:       req.Title,
		Description: req.Description,
		Department:  models.Department(req.Department),
		CreatorID:   creatorID,
		Deadline:    req.Deadline,
		AllowLate:   req.AllowLate,
		Version:     0,
	}

	if err := s.homeworkDao.Create(homework); err != nil {
		return nil, errcode.ServerError, err
	}

	// ========== 极简并发示例 ==========
	go func() {
		// 查询同部门学员数量
		var count int64
		configs.DB.Model(&models.User{}).Where("department = ? AND role = ?", req.Department, "student").Count(&count)

		// 创建channel
		resultChan := make(chan string, 1)

		// 启动goroutine模拟通知
		go func() {
			time.Sleep(1 * time.Second) // 模拟耗时操作
			resultChan <- fmt.Sprintf("已通知 %d 名学员", count)
		}()

		// 使用select处理超时
		select {
		case msg := <-resultChan:
			log.Println("✅", msg)
		case <-time.After(2 * time.Second):
			log.Println("⚠️ 通知超时")
		}
	}()
	// ========== 极简并发示例结束 ==========

	// 重新查询以获取关联的Creator信息
	homework, _ = s.homeworkDao.FindByID(homework.ID)
	return homework, errcode.Success, nil
}

// 获取作业详情
func (s *HomeworkService) GetByID(id uint) (*models.Homework, int, error) {
	homework, err := s.homeworkDao.FindByID(id)
	if err != nil {
		return nil, errcode.HomeworkNotFound, err
	}
	return homework, errcode.Success, nil
}

// 获取作业列表
func (s *HomeworkService) List(department string, page, pageSize int) (*HomeworkListResponse, int, error) {
	// 参数校验
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	// 查询数据
	homeworks, total, err := s.homeworkDao.List(department, page, pageSize)
	if err != nil {
		return nil, errcode.ServerError, err
	}

	// 构造响应
	items := make([]HomeworkItemResponse, 0, len(homeworks))
	for _, hw := range homeworks {
		items = append(items, HomeworkItemResponse{
			ID:              hw.ID,
			Title:           hw.Title,
			Description:     hw.Description,
			Department:      string(hw.Department),
			DepartmentLabel: models.DepartmentLabel[hw.Department],
			Creator: Creator{
				ID:              hw.Creator.ID,
				Nickname:        hw.Creator.Nickname,
				Department:      string(hw.Creator.Department),
				DepartmentLabel: models.DepartmentLabel[hw.Creator.Department],
			},
			Deadline:  hw.Deadline,
			AllowLate: hw.AllowLate,
			Version:   hw.Version,
			CreatedAt: hw.CreatedAt,
		})
	}

	return &HomeworkListResponse{
		Total: total,
		Items: items,
		Page:  page,
		Size:  pageSize,
	}, errcode.Success, nil
}

// 更新作业（带乐观锁）
func (s *HomeworkService) Update(operatorID uint, req *UpdateHomeworkRequest) (*models.Homework, int, error) {
	// 获取原作业
	homework, err := s.homeworkDao.FindByID(req.ID)
	if err != nil {
		return nil, errcode.HomeworkNotFound, err
	}
	log.Printf("🔍 截止时间验证: deadline=%v, now=%v", req.Deadline, time.Now())
	if req.Deadline.Before(time.Now()) {
		log.Printf("❌ 截止时间已过: %v < %v", req.Deadline, time.Now())
		return nil, errcode.InvalidParam, models.ErrDeadlinePassed
	}
	log.Printf("✅ 截止时间验证通过")
	log.Printf("🔍 版本号验证: req.Version=%d, db.Version=%d", req.Version, homework.Version)
	if homework.Version != req.Version {
		log.Printf("❌ 版本号不匹配")
		return nil, errcode.HomeworkVersionExpired, nil
	}
	log.Printf("✅ 版本号验证通过")
	// 验证操作者权限
	operator, err := s.userDao.FindByID(operatorID)
	if err != nil {
		return nil, errcode.UserNotFound, err
	}
	if !homework.CanModify(operator) {
		return nil, errcode.HomeworkNotPermission, nil
	}

	// 验证截止时间
	if req.Deadline.Before(time.Now()) {
		return nil, errcode.InvalidParam, models.ErrDeadlinePassed
	}

	// 验证版本号
	if homework.Version != req.Version {
		return nil, errcode.HomeworkVersionExpired, nil
	}

	// 更新字段
	homework.Title = req.Title
	homework.Description = req.Description
	homework.Deadline = req.Deadline
	homework.AllowLate = req.AllowLate

	// 保存
	if err := s.homeworkDao.Update(homework); err != nil {
		if err == dao.ErrVersionExpired {
			return nil, errcode.HomeworkVersionExpired, err
		}
		return nil, errcode.ServerError, err
	}

	homework, _ = s.homeworkDao.FindByID(homework.ID)
	return homework, errcode.Success, nil
}

// 删除作业
func (s *HomeworkService) Delete(operatorID, homeworkID uint) (int, error) {
	// 获取作业
	homework, err := s.homeworkDao.FindByID(homeworkID)
	if err != nil {
		return errcode.HomeworkNotFound, err
	}

	// 验证操作者权限
	operator, err := s.userDao.FindByID(operatorID)
	if err != nil {
		return errcode.UserNotFound, err
	}
	if !homework.CanModify(operator) {
		return errcode.HomeworkNotPermission, nil
	}

	// 删除
	if err := s.homeworkDao.Delete(homeworkID); err != nil {
		return errcode.ServerError, err
	}

	return errcode.Success, nil
}
