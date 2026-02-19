package configs

import (
	"homework-system/models" //导入模型包，用于自动建表
	"log"                    //日志
	"os"                     //读取环境的变化量
	"time"                   //时间处理

	"gorm.io/driver/mysql" //MySQL数据库的驱动
	"gorm.io/gorm"         //ORM框架
	"gorm.io/gorm/logger"  //GORM日志
)

var DB *gorm.DB //全局数据库连接对象，其他包可以通过configs.DB来操作数据库

type Config struct {
	DBHost     string //数据库主机地址
	DBPort     string //数据库的端口
	DBUser     string //数据库用户名
	DBPassword string //数据库密码
	DBName     string //数据库名
	Port       string //服务器端口
}

func LoadConfig() *Config { //加载配置，优先从环境变量读取，没有再用默认值
	// 从环境变量读取，没有则用默认值
	return &Config{
		DBHost:     getEnv("DB_HOST", "127.0.0.1"),       //默认本地
		DBPort:     getEnv("DB_PORT", "3306"),            //默认是3306
		DBUser:     getEnv("DB_USER", "root"),            //默认是root
		DBPassword: getEnv("DB_PASSWORD", "520405"),      //默认密码
		DBName:     getEnv("DB_NAME", "homework_system"), //默认数据库
		Port:       getEnv("PORT", "8080"),               //默认8080端口
	}
}

func getEnv(key, defaultValue string) string { //读取环境变量，如果没有则返回默认值
	if value := os.Getenv(key); value != "" {  //如果环境变量存在，返回环境变量
		return value
	}
	return defaultValue //否则返回默认值
}

func Init() { //初始化函数
	config := LoadConfig() //加载配置
	initDB(config)         //初始化数据库
}

func initDB(config *Config) { //初始化数据库
	// 构建 DSN 数据库连接字符串，格式为：用户名：密码@tcp（主机：端口）/数据库名？参数
	dsn := config.DBUser + ":" + config.DBPassword + "@tcp(" + config.DBHost + ":" + config.DBPort + ")/" + config.DBName + "?charset=utf8mb4&parseTime=True&loc=Local"
	//charset=utf8mb4.支持表情符号
	//parseTime=True.将数据库时间自动解析为go的time.Time包
	//loc=Local.使用本地时区
	gormConfig := &gorm.Config{ // GORM日志配置
		Logger: logger.Default.LogMode(logger.Info), //日志级别为info，会打印SQL语句
	}

	db, err := gorm.Open(mysql.Open(dsn), gormConfig) //连接数据库
	if err != nil {
		log.Fatal("❌ 数据库连接失败: ", err) //失败就退出程序
	}

	err = db.AutoMigrate( // 自动建表，根据models（模型层）定义的结构体，自动创建对应的数据库表
		&models.User{},       //创建User表
		&models.Homework{},   //创建Homework表
		&models.Submission{}, //创建Submission表
	)
	if err != nil {
		log.Fatal("❌ 数据库迁移失败: ", err)
	}

	// 设置连接池
	sqlDB, err := db.DB() //获取底层的sql.DB的对象
	if err != nil {
		log.Fatal("❌ 获取数据库连接失败: ", err)
	}
	sqlDB.SetMaxIdleConns(10)           //设置最大空闲连接数
	sqlDB.SetMaxOpenConns(100)          //设置最大打开连接数
	sqlDB.SetConnMaxLifetime(time.Hour) //设置连接最大生命周期数

	DB = db //赋值给全局变量，供其他层的包使用
	log.Println("✅ 数据库连接成功")
}
