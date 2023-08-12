package config

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
)

// Mysql：数据库配置
type Mysql struct {
	Host      string
	Port      int
	Database  string
	Username  string
	Password  string
	Charset   string //utf8mb4
	ParseTime bool   `toml:"parse_time"` //解析mysql中的时间类型到go中的time类型
	Loc       string //时区
}

// Redis：缓存配置
type Redis struct {
	IP       string
	Port     int
	Database int
}

// Server：服务器配置
type Server struct {
	IP   string
	Port int
}

// Path：路径配置
type Path struct {
	FfmpegPath       string `toml:"ffmpeg_path"`        //ffmpeg的路径
	StaticSourcePath string `toml:"static_source_path"` //静态资源路径
}

// Config：配置
type Config struct {
	DB     Mysql `toml:"mysql"`
	RDB    Redis `toml:"redis"`
	Server `toml:"server"`
	Path   `toml:"path"`
}

// Global：全局配置
var Global Config

// ensurePathValid：确保路径的有效性，检查静态资源路径和ffmpeg路径是否存在，并进行相应的处理
func ensurePathValid() {
	var err error
	if _, err = os.Stat(Global.StaticSourcePath); os.IsNotExist(err) {
		if err = os.Mkdir(Global.StaticSourcePath, 0755); err != nil {
			log.Fatalf("mkdir error:path %s", Global.StaticSourcePath)
		}
	}
	if _, err = os.Stat(Global.FfmpegPath); os.IsNotExist(err) {
		if _, err = exec.Command("ffmpeg", "-version").Output(); err != nil {
			log.Fatalf("ffmpeg not valid %s", Global.FfmpegPath)
		} else {
			Global.FfmpegPath = "ffmpeg"
		}
	} else {
		Global.FfmpegPath, err = filepath.Abs(Global.FfmpegPath)
		if err != nil {
			log.Fatalln("get abs path failed:", Global.FfmpegPath)
		}
	}
	//把资源路径转化为绝对路径，防止调用ffmpeg命令失效
	Global.StaticSourcePath, err = filepath.Abs(Global.StaticSourcePath)
	if err != nil {
		log.Fatalln("get abs path failed:", Global.StaticSourcePath)
	}
}

// init：在包初始化时被调用，加载配置文件并初始化，如解码配置文件、修剪字符串和验证路径
func init() {
	if _, err := toml.DecodeFile("./config/config.toml", &Global); err != nil {
		panic(err)
	}
	//去除左右的空格
	strings.Trim(Global.Server.IP, " ")
	strings.Trim(Global.RDB.IP, " ")
	strings.Trim(Global.DB.Host, " ")
	//保证路径正常
	ensurePathValid()
}

// DBConnectString：填充得到数据库连接字符串
func DBConnectString() string {
	arg := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=%s&parseTime=%v&loc=%s",
		Global.DB.Username, Global.DB.Password, Global.DB.Host, Global.DB.Port, Global.DB.Database,
		Global.DB.Charset, Global.DB.ParseTime, Global.DB.Loc)
	log.Println(arg)
	return arg
}
