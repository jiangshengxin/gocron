package app

import (
	"os"
	"scheduler/modules/crontask"
	"scheduler/models"
	"runtime"
	"scheduler/modules/utils"
)

var  (
	AppDir string    // 应用根目录
	ConfDir string   // 配置目录
	LogDir string    // 日志目录
	DataDir string   // 数据目录，存放session文件等
	AppConfig string // 应用配置文件
	Installed bool   // 应用是否安装过
	CronTask crontask.CronTask // 定时任务
)

func init() {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	AppDir = wd
	ConfDir = AppDir + "/conf"
	LogDir  = AppDir + "/log"
	DataDir = AppDir + "/data"
	AppConfig = AppDir + "/app.ini"
	checkDirExists(ConfDir, LogDir, DataDir)
	// ansible配置文件目录
	os.Setenv("ANSIBLE_CONFIG", ConfDir)
	Installed = IsInstalled()
	if Installed {
		initResource()
	}
}

// 判断应用是否安装过
func IsInstalled() bool {
	_, err := os.Stat(ConfDir + "/install.lock")
	if os.IsNotExist(err) {
		return false
	}

	return true
}

// 检测环境
func CheckEnv()  {
	// ansible不支持安装在windows上, windows只能作为被控机
	if runtime.GOOS == "windows" {
		panic("不支持在windows上运行")
	}
	_, err := utils.ExecShell("ansible", "--version")
	if err != nil {
		panic(err)
	}
	_, err = utils.ExecShell("ansible-playbook", "--version")
	if err != nil {
		panic("ansible-playbook not found")
	}
}

// 创建安装锁文件
func CreateInstallLock() error {
	_, err := os.Create(ConfDir + "/install.lock")
	if err != nil {
		utils.RecordLog("创建安装锁文件失败")
	}

	return err
}


// 初始化资源
func initResource()  {
	crontask.DefaultCronTask = crontask.CreateCronTask()

	models.Db = models.CreateDb(AppConfig)
}

// 检测目录是否存在
func checkDirExists(path... string)  {
	for _, value := range(path) {
		_, err := os.Stat(value)
		if os.IsNotExist(err) {
			panic(value + "目录不存在")
		}
		if os.IsPermission(err) {
			panic(value + "目录无权限操作")
		}
	}
}