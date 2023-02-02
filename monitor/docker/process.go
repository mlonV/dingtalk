package docker

import (
	"context"
	"fmt"
	"regexp"

	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/docker/docker/client"
	"github.com/mlonV/dingtalk/config"
	"github.com/mlonV/dingtalk/prome"
	"github.com/mlonV/dingtalk/utils"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	dm config.MonitorDocker

	// 全局的sync.Map 放dockerID和ProcessStatus结构体
	dMap   sync.Map
	ticker *time.Ticker
	psChan chan ProcessStatus // 传入进程结构体

)

func init() {

	//  DockerMonitor
	dm = config.Conf.MonitorDocker
	psChan = make(chan ProcessStatus, 100)

	fmt.Printf("%#v\n", dm)
}

// docker内进程
type ProcessStatus struct {
	dockerCli     *client.Client // dockerCli
	Addr          string         // 所属的主机
	ContainerName string         // 容器名
	ContainerID   string         // 容器ID
	ProcessName   string         // ProcessName 进程名
	Alive         float64        // 初始值-1  0启动 1down

	// PID相关
	PID      string // 记录PID, pid变化告警（说明进程重启了
	IsChange bool   //是否变化过了
	Counter  int64  // 计数器
	// PromePIDGauge

	// prometheus Gauge
	PromeGauge    prometheus.Gauge
	PromePIDGauge prometheus.Gauge

	Message string // Msg 发送message
}

// 定时器
func Ticker() {

	// 启动处理chan进程
	go HandleChan(psChan)
	// interval := time.Duration(esalarm.Time) * time.Minute
	interval := time.Second * time.Duration(dm.Interval)
	ticker = time.NewTicker(interval)
	config.Log.Info("Start Monitor Process/PID ,Interval: %ds ", dm.Interval)
	for {
		// 调用Reset方法对timer对象进行定时器重置
		// 	ticker.Reset(interval)
		<-ticker.C
		Worker()
	}
	// ticker.Stop()
}

// 连到服务器上使用docker api获取到Running的容器名和ID
func Worker() {
	// 获取每台服务器上容器名和id对应的Map
	for _, host := range dm.Hosts {
		// 获取dockercli
		dockerCli, err := utils.NewDockerCli(dm.Username, host, fmt.Sprint(dm.Port))
		config.Log.Debug("dockerCli.ClientVersion() : %s ", dockerCli.ClientVersion())
		if err != nil {
			config.Log.Error(err.Error())
		}
		// 获取到容器id和ProcessStatus的对应的map
		// 先从dMap里面拿，拿不到传个新的进去
		err = GetRunningDockerInfo(dockerCli, host)
		if err != nil {
			config.Log.Error(err.Error())
		}
	}

}

// 传入主机信息，获取到running状态的docker id/name/ 容器内进程号/进程名
func GetRunningDockerInfo(dockerCli *client.Client, addr string) error {
	// tm的忘记关闭连接了
	defer dockerCli.Close()
	cList, err := utils.GetContainerByDocker(dockerCli)

	if err != nil {
		return err
	}
	// 遍历每一个容器
	info, err := dockerCli.Info(context.Background())
	if err != nil {
		return err
	}

	fmt.Printf("%s  || %#v", dockerCli.DaemonHost(), info.Name)
	for _, container := range cList {
		// 只搞运行中的容器
		if container.State == "running" {
			containerName := strings.Split(container.Names[0], "/")[1]
			value, ok := dMap.Load(containerName)
			var ps ProcessStatus
			if !ok {
				ps.ContainerName = containerName
				ps.ContainerID = container.ID
				ps.Addr = addr
				ps.Alive = -1
				// 从dockerCli来获取主机名
				info, err := dockerCli.Info(context.Background())
				if err != nil {
					config.Log.Error("Get Docker Info err: %s", err)
				}
				// 如果没有prometheus-cli 则分配一个
				if ps.PromeGauge == nil {
					ps.PromeGauge = prometheus.NewGauge(prometheus.GaugeOpts{
						Namespace: "victory",
						Subsystem: "go",
						Name:      "alive",
						Help:      "容器内进程存活指标 is Alive",
						ConstLabels: map[string]string{
							"app":      "game",
							"hostname": info.Name,
							"address":  addr,
							"app_name": ps.ContainerName,
						},
					})
				}

				if ps.PromePIDGauge == nil {
					ps.PromePIDGauge = prometheus.NewGauge(prometheus.GaugeOpts{
						Namespace: "victory",
						Subsystem: "go",
						Name:      "pid",
						Help:      "容器内进程存活指标 is Alive",
						ConstLabels: map[string]string{
							"app":      "game",
							"hostname": info.Name,
							"address":  addr,
							"app_name": ps.ContainerName,
						},
					})
				}
			} else {
				ps, _ = (value).(ProcessStatus)
			}

			ps.dockerCli = dockerCli

			// 从全局的dMap中查是否已经有了容器进程的键值信息
			// 查出来进程号和进程
			// 先判断是不是gamex的容器（gamex的容器名都是gamex开头），如果是gamex则按照gamex处理
			process := dm.Process
			if IsGameX(containerName) {
				config.Log.Info("当前ContainerName: %s处理gamex的获取流程,设置dm.Process = %s", containerName, dm.GameXPath)
				process = dm.GameXPath
			}
			pidCmd := []string{"bash", "-c", fmt.Sprintf(`ps -ef|grep "%s" |grep -v grep|awk '{print $2}'`, process)}
			pid, err := utils.ExecCmd(ps.dockerCli, pidCmd, container.ID)
			if err != nil {
				return err
			}
			psnameCmd := []string{"bash", "-c", fmt.Sprintf(`ps -ef|grep "%s" |grep -v grep|awk '{print $8}'`, process)}
			psname, err := utils.ExecCmd(ps.dockerCli, psnameCmd, container.ID)
			if err != nil {
				return err
			}
			pid = strings.Split(pid, "\n")[0]
			psname = strings.Split(psname, "\n")[0]

			config.Log.Debug("当前隐藏字符集 %#v,%#v,process: %s", pid, psname, process)
			config.Log.Info("当前ContainerName [%s], PID : [%s], ProcessName : [%s]", ps.ContainerName, pid, psname)

			// 查不出来进程挂掉,Alive !=-1则是已经添加到Map的情况
			if ps.Alive == 0 && (pid == "" || psname == "") {
				ps.Alive = 1
				ps.IsChange = true
				ps.PromeGauge.Set(ps.Alive)
				config.Log.Info("PromeGauge Set [%s] Value : %f", ps.ContainerName, ps.Alive)
				ps.Message = fmt.Sprintf("Host: [%s] ,Container: [%s] is Down", ps.Addr, ps.ContainerName)
				psChan <- ps
				// 跳出之前先把值存进dMap
				dMap.Store(ps.ContainerName, ps)
				continue
			}

			// 延迟告警
			if ps.Alive == 1 && (pid == "" || psname == "") {
				ps.Alive = 1
				ps.IsChange = true
				ps.PromeGauge.Set(ps.Alive)
				config.Log.Info("PromeGauge Set [%s] Value : %f", ps.ContainerName, ps.Alive)
				ps.Counter += 1
				if ps.Counter > dm.Num {
					ps.Counter = 0
					ps.Message = fmt.Sprintf("Host: [%s] ,Container: [%s] is Down", ps.Addr, ps.ContainerName)
					psChan <- ps
					dMap.Store(ps.ContainerName, ps)
					continue
				}
				dMap.Store(ps.ContainerName, ps)
			}
			// pid没变化则继续
			if pid == ps.PID && !ps.IsChange {
				continue
			}
			if ps.PID != "" && pid != ps.PID {
				ps.IsChange = true
				dMap.Store(ps.ContainerName, ps)

			}

			if pid == "" {
				continue
			}

			// 查到结束前查看pid是否和之前的一样
			if ps.IsChange {
				ps.Counter += 1
				if ps.Counter > dm.Num {
					ps.IsChange = false
					ps.Counter = 0
					ps.Message = fmt.Sprintf("Host: [%s] ,Container: [%s] Pid Change Resloved", ps.Addr, ps.ContainerName)
					psChan <- ps
				}
				ps.PromePIDGauge.Set(GetGaugeValue(ps.Counter, pid))
				config.Log.Info("PromePIDGauge Set [%s] Value : %f", ps.ContainerName, GetGaugeValue(ps.Counter, pid))
				dMap.Store(ps.ContainerName, ps)
			}
			ps.Alive = 0
			ps.PromeGauge.Set(ps.Alive)
			config.Log.Info("PromeGauge Set [%s] Value : %f", ps.ContainerName, ps.Alive)
			ps.PromePIDGauge.Set(GetGaugeValue(ps.Counter, pid))
			config.Log.Info("PromePIDGauge Set [%s] Value : %f", ps.ContainerName, GetGaugeValue(ps.Counter, pid))
			ps.PID = pid
			ps.ProcessName = psname
			// 发送注册 prome
			config.Log.Info("注册PromeGauge到prometheus Register : %s", ps.ContainerName)
			config.Log.Info("注册PromePIDGauge到prometheus Register : %s", ps.ContainerName)
			psChan <- ps
			dMap.Store(ps.ContainerName, ps)
		}

	}
	return nil
}

// 判断容器名是不是gamex开头的
// 传入容器名，返回true/false
func IsGameX(cname string) bool {
	reg, err := regexp.Compile(`gamex`)
	if err != nil {
		config.Log.Error("IsGamex Compile err : ", err.Error())
	}
	s := reg.FindAllString(cname, 10)
	if s != nil {
		config.Log.Info("%s is Gamex ", cname)
		config.Log.Info("%#v", s)
		return true
	}
	return false
}

// 处理chan
func HandleChan(psChan chan ProcessStatus) {
	for {
		ps := <-psChan

		prome.PromeRegister.Register(ps.PromeGauge)
		prome.PromeRegister.Register(ps.PromePIDGauge)
		if ps.Alive == 0 {
			// 启动通知
			config.Log.Warning("process up by HandleChan : %#v", ps)
		}
		if ps.Alive == 1 {
			// 挂掉告警
			config.Log.Warning("process down by HandleChan : %#v", ps)
		}
	}
}

// 获取用来设置promeSet的float64值
func GetGaugeValue(counter int64, pid string) float64 {
	gaugeValue, _ := strconv.ParseFloat(fmt.Sprintf("%d.%s", counter, pid), 64)
	return gaugeValue
}
