package server

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"microservice/util"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

var (
	microserviceConf = &MicroserviceConf{
		Port: 8080,
		Prometheus: PrometheusConf{
			SwitchOn: true,
			Port:     8081,
		},
		ServiceName: "koala_server",
		Register: RegisterConf{
			SwitchOn: false,
		},
		Log: LogConf{
			Level: "debug",
			Dir:   "./logs/",
		},
	}
)

type MicroserviceConf struct {
	Port        int            `yaml:"port"`
	Prometheus  PrometheusConf `yaml:"prometheus"`
	Register    RegisterConf   `yaml:"register"`
	ServiceName string         `yaml:"service_name"`
	Log         LogConf        `yaml:"log"`
	Limit       LimitConf      `yaml:"limit"`
	Trace       TraceConf      `yaml:"trace"`

	//内部的配置项
	ConfigDir  string `yaml:"-"`
	RootDir    string `yaml:"-"`
	ConfigFile string `yaml:"-"`
}

type LimitConf struct {
	QPSLimit int  `yaml:"qps"`
	SwitchOn bool `yaml:"switch_on"`
}

type RegisterConf struct {
	SwitchOn     bool          `yaml:"switch_on"`
	RegisterPath string        `yaml:"register_path"`
	Timeout      time.Duration `yaml:"timeout"`
	HeartBeat    int64         `yaml:"heart_beat"`
	RegisterName string        `yaml:"register_name"`
	RegisterAddr string        `yaml:"register_addr"`
}

type PrometheusConf struct {
	SwitchOn bool `yaml:"switch_on"`
	Port     int  `yaml:"port"`
}

type LogConf struct {
	Level      string `yaml:"level"`
	Dir        string `yaml:"path"`
	ChanSize   int    `yaml:"chan_size"`
	ConsoleLog bool   `yaml:"console_log"`
}

func initDir(serviceName string) (err error) {
	exeFilePath, err := filepath.Abs(os.Args[0])
	if err != nil {
		return
	}

	if runtime.GOOS == "windows" {
		exeFilePath = strings.Replace(exeFilePath, "\\", "/", -1)
	}

	lastIdx := strings.LastIndex(exeFilePath, "/")
	if lastIdx < 0 {
		err = fmt.Errorf("invalid exe path:%v", exeFilePath)
		return
	}
	//C:/project/src/xxx/
	microserviceConf.RootDir = path.Join(strings.ToLower(exeFilePath[0:lastIdx]), "..")
	microserviceConf.ConfigDir = path.Join(microserviceConf.RootDir, "./conf/", util.GetEnv())
	microserviceConf.ConfigFile = path.Join(microserviceConf.ConfigDir, fmt.Sprintf("%s.yaml", serviceName))
	return
}

func InitConfig(serviceName string) (err error) {

	err = initDir(serviceName)
	if err != nil {
		return
	}

	data, err := ioutil.ReadFile(microserviceConf.ConfigFile)
	if err != nil {
		return
	}

	err = yaml.Unmarshal(data, &microserviceConf)
	if err != nil {
		return
	}

	fmt.Printf("init koala conf succ, conf:%#v\n", microserviceConf)
	return
}

type TraceConf struct {
	SwitchOn   bool    `yaml:"switch_on"`
	ReportAddr string  `yaml:"report_addr"`
	SampleType string  `yaml:"sample_type"`
	SampleRate float64 `yaml:"sample_rate"`
}

func GetConfigDir() string {
	return microserviceConf.ConfigDir
}

func GetRootDir() string {
	return microserviceConf.RootDir
}

func GetServerPort() int {
	return microserviceConf.Port
}

func GetConf() *MicroserviceConf {
	return microserviceConf
}
