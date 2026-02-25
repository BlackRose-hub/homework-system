🚀 红岩作业管理系统 - BlackRose的寒假作业



> 一个让老登轻松布置、小登快乐交作业的Go后端项目



---



&nbsp;📋 项目简介



\*\*一句话说明\*\*： 一个面向多部门的作业系统



\*\*解决的问题\*\*：

\- 老登（管理员）想发作业，但不想一个个通知

\- 小登（学员）想交作业，但不想搞混部门

\- 各部门作业要分开管理，不能乱



\*\*最终效果\*\*：7个部门、2种角色、15个接口，完整跑通！







🛠️ 技术栈



| 技术 | 版本 | 用来干嘛 |

|------|------|----------|

| Go | 1.21 | 写后端逻辑 |

| Gin | v1.9.1 | 处理HTTP请求 |

| MySQL | 8.0 | 存数据 |

| GORM | v1.25.5 | 不用写SQL |

| JWT | v5 | 双token认证 |

| bcrypt | - | 密码加密 |



\*\*特别说明\*\*：GORM + MySQL 真香，不用写SQL太爽了！



---



📁 项目结构

homework-system/

├── cmd/ # 程序入口

│ └── main.go

├── configs/ # 配置文件

│ └── config.go # 数据库连接配置

├── dao/ # 数据访问层

│ ├── user.go

│ ├── homework.go

│ └── submission.go

├── handler/ # HTTP处理器

│ ├── user.go

│ ├── homework.go

│ └── submission.go

├── middleware/ # 中间件

│ ├── auth.go # JWT认证

│ └── permission.go # 权限控制

├── models/ # 数据模型

│ ├── user.go

│ ├── homework.go

│ └── submission.go

├── pkg/ # 工具包

│ ├── jwt/ # JWT工具

│ ├── response/ # 统一响应

│ └── errcode/ # 错误码定义

├── router/ # 路由定义

│ └── router.go

├── service/ # 业务逻辑层

│ ├── user.go

│ ├── homework.go

│ └── submission.go

└── README.md



---



\## 4. 已实现功能清单



\### 👤 用户模块

\- 用户注册（支持部门选择）

\- 用户登录（JWT双Token认证）

\- 刷新Token（Refresh Token机制）

\- 获取用户信息（带部门标签）

\- 注销账号（软删除）



\### 📝 作业模块

\- 发布作业（管理员权限）

\- 作业列表（按部门筛选、分页查询）

\- 作业详情（带发布者信息和部门标签）

\- 修改作业（同部门管理员 + 乐观锁并发控制）

\- 删除作业（同部门管理员 + 软删除）



\### 📤 提交模块

\- 提交作业（自动判断是否迟交）

\- 我的提交（学生查看自己的提交和评语）

\- 部门提交（管理员查看本部门学员提交）

\- 批改作业（打分、写评语）

\- 标记优秀作业

\- 优秀作业展示（公开接口，无需登录）



---



5\. 项目亮点



&nbsp;✨ JWT双Token认证

\- Access Token（15分钟过期）：用于接口认证

\- Refresh Token（7天过期）：用于无感刷新

\- 兼顾安全性与用户体验



&nbsp;✨ 乐观锁并发控制

\- 解决多个管理员同时修改作业的冲突问题

\- 数据库添加Version字段，更新时校验版本号

\- 防止数据被覆盖



&nbsp;✨ Go并发编程示例

\- 发布作业时异步通知学员

\- 使用goroutine + channel + select

\- 不阻塞主流程，提升响应速度



&nbsp;✨ 部门枚举+中文标签

\- 请求时使用枚举值（backend/frontend等）

\- 响应时同时返回枚举值和中文标签（后端/前端等）

\- 保证数据规范性和前端友好性



---



6\. API 接口说明



&nbsp; 6.1 通用规范



\*\*请求头\*\*

&nbsp;Content-Type: application/json

Authorization: Bearer <access\_token> # 需要认证的接口必填



\*\*统一响应格式\*\*

```json

{

&nbsp; "code": 1000,

&nbsp; "message": "成功",

&nbsp; "data": {}

}

错误码示例

{

&nbsp; "code": 2001,

&nbsp; "message": "用户不存在",

&nbsp; "data": null

}

6.2 部门枚举值

枚举值	显示标签	说明

backend	后端	后端开发

frontend	前端	前端开发

sre	SRE	运维工程

product	产品	产品设计

design	视觉设计	UI/UX设计

android	Android	Android开发

ios	iOS	iOS开发

6.3 接口列表

用户模块

方法	路径	功能	认证

POST	/api/auth/register	用户注册	否

POST	/api/auth/login	用户登录	否

POST	/api/auth/refresh	刷新Token	否

GET	/api/user/profile	获取个人信息	是

DELETE	/api/user/account	注销账号	是

作业模块

方法	路径	功能	权限

POST	/api/homeworks	发布作业	管理员

GET	/api/homeworks	作业列表	登录用户

GET	/api/homeworks/:id	作业详情	登录用户

PUT	/api/homeworks/:id	修改作业	管理员+同部门

DELETE	/api/homeworks/:id	删除作业	管理员+同部门

提交模块

方法	路径	功能	权限

POST	/api/submissions	提交作业	学生

GET	/api/submissions/my	我的提交	学生

GET	/api/submissions/department	部门提交	管理员

POST	/api/submissions/review	批改作业	管理员

GET	/api/submissions/excellent	优秀作业	公开

7\. 本地运行指南

环境要求

Go 1.21+

MySQL 8.0+

Git

8\. 作者信息

姓名：罗雅岚



完成时间：2026年2月20日



GitHub：https://github.com/BlackRose-hub/homework-system

