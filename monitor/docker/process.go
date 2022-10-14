package docker

import (
	"fmt"
	"log"
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

type ProcessStatus struct {
	dockerCli     *client.Client // dockerCli
	Host          string         // 所属的主机
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
	log.Println("开始遍历所有的Host 的 容器 ")
	for _, host := range dm.Hosts {
		// 获取dockercli
		dockerCli, err := utils.NewDockerCli(dm.Username, host, fmt.Sprint(dm.Port))
		dockerCli.ClientVersion()
		if err != nil {
			log.Println(err.Error())
		}
		// 获取到容器id和ProcessStatus的对应的map
		// 先从dMap里面拿，拿不到传个新的进去
		GetRunningDockerInfo(dockerCli, host)
		if err != nil {
			log.Println(err.Error())
		}
	}

}

// 传入主机信息，获取到running状态的docker id/name/ 容器内进程号/进程名
func GetRunningDockerInfo(dockerCli *client.Client, host string) error {
	defer dockerCli.Close()

	cList, err := utils.GetContainerByDocker(dockerCli)
	// tm的忘记关闭连接了

	if err != nil {
		return err
	}
	// 遍历每一个容器

	for _, container := range cList {
		// 只搞运行中的容器
		if container.State == "running" {
			containerName := strings.Split(container.Names[0], "/")[1]
			value, ok := dMap.Load(containerName)
			var ps ProcessStatus
			// log.Printf("dMap : %#v,%t,%s", dMap, ok, containerName)
			if !ok {
				ps.ContainerName = containerName
				ps.ContainerID = container.ID
				ps.Host = host
				ps.dockerCli = dockerCli
				ps.Alive = -1
				// 如果没有prometheus-cli 则分配一个
				if ps.PromeGauge == nil {
					ps.PromeGauge = prometheus.NewGauge(prometheus.GaugeOpts{
						Namespace: "victory",
						Subsystem: "go",
						Name:      ps.ContainerName,
						Help:      "容器内进程存活指标 is Alive",
					})
				}

				if ps.PromePIDGauge == nil {
					ps.PromePIDGauge = prometheus.NewGauge(prometheus.GaugeOpts{
						Namespace: "victory",
						Subsystem: "go_pid",
						Name:      ps.ContainerName,
						Help:      "容器内进程存活指标 is Alive",
					})
				}
			} else {
				ps, _ = (value).(ProcessStatus)
				ps.dockerCli = dockerCli
			}

			// 从全局的dMap中查是否已经有了容器进程的键值信息
			// 查出来进程号和进程
			pidCmd := []string{"bash", "-c", fmt.Sprintf("ps -ef|grep '%s' |grep -v grep|awk '{print $2}'", dm.Process)}
			pid, err := utils.ExecCmd(ps.dockerCli, pidCmd, container.ID)
			if err != nil {
				return err
			}
			psnameCmd := []string{"bash", "-c", fmt.Sprintf("ps -ef|grep '%s' |grep -v grep|awk '{print $8}'", dm.Process)}
			psname, err := utils.ExecCmd(ps.dockerCli, psnameCmd, container.ID)
			if err != nil {
				return err
			}
			//  docker 容器内执行命令会带有隐藏字符ascii字符，用正则替换掉
			reg, err := regexp.Compile(`[\x00-\x1F]`)

			if err != nil {
				return err
			}
			bPID := reg.ReplaceAll([]byte(strings.Replace(pid, "\n", "", -1)), []byte{})
			bpsname := reg.ReplaceAll([]byte(strings.Replace(psname, "\n", "", -1)), []byte{})
			log.Printf("ContainerName [%s], PID : [%s], ProcessName : [%s]", ps.ContainerName, string(bPID), string(bpsname))

			// 查不出来进程挂掉,Alive !=-1则是已经添加到Map的情况
			if ps.Alive == 0 && (string(bPID) == "" || string(bpsname) == "") {
				ps.Alive = 1
				ps.IsChange = true
				ps.PromeGauge.Set(ps.Alive)
				log.Printf("PromeGauge Set [%s] Value : %f", ps.ContainerName, ps.Alive)
				ps.Message = fmt.Sprintf("Host: [%s] ,Container: [%s] is Down", ps.Host, ps.ContainerName)
				psChan <- ps
				// 跳出之前先把值存进dMap
				dMap.Store(ps.ContainerName, ps)
				continue
			}

			// 延迟告警
			if ps.Alive == 1 && (string(bPID) == "" || string(bpsname) == "") {
				ps.Alive = 1
				ps.IsChange = true
				ps.PromeGauge.Set(ps.Alive)
				log.Printf("PromeGauge Set [%s] Value : %f", ps.ContainerName, ps.Alive)
				ps.Counter += 1
				if ps.Counter > dm.Num {
					ps.Counter = 0
					ps.Message = fmt.Sprintf("Host: [%s] ,Container: [%s] is Down", ps.Host, ps.ContainerName)
					psChan <- ps
					dMap.Store(ps.ContainerName, ps)
					continue
				}
				dMap.Store(ps.ContainerName, ps)
			}
			// pid没变化则继续
			if string(bPID) == ps.PID && !ps.IsChange {
				continue
			}
			if ps.PID != "" && string(bPID) != ps.PID {
				ps.IsChange = true
				dMap.Store(ps.ContainerName, ps)

			}

			if string(bPID) == "" {
				continue
			}

			// 查到结束前查看pid是否和之前的一样
			if ps.IsChange {
				ps.Counter += 1
				if ps.Counter > dm.Num {
					ps.IsChange = false
					ps.Counter = 0
					ps.Message = fmt.Sprintf("Host: [%s] ,Container: [%s] Pid Change Resloved", ps.Host, ps.ContainerName)
					psChan <- ps
				}
				ps.PromePIDGauge.Set(GetGaugeValue(ps.Counter, string(bPID)))
				log.Printf("PromePIDGauge Set [%s] Value : %f", ps.ContainerName, GetGaugeValue(ps.Counter, string(bPID)))
				dMap.Store(ps.ContainerName, ps)
			}
			ps.Alive = 0
			ps.PromeGauge.Set(ps.Alive)
			log.Printf("PromeGauge Set [%s] Value : %f", ps.ContainerName, ps.Alive)
			ps.PromePIDGauge.Set(GetGaugeValue(ps.Counter, string(bPID)))
			log.Printf("PromePIDGauge Set [%s] Value : %f", ps.ContainerName, GetGaugeValue(ps.Counter, string(bPID)))
			ps.PID = string(bPID)
			ps.ProcessName = string(bpsname)
			// 发送注册 prome
			log.Printf("注册PromeGauge到prometheus Register : %s", ps.ContainerName)
			log.Printf("注册PromePIDGauge到prometheus Register : %s", ps.ContainerName)
			psChan <- ps
			dMap.Store(ps.ContainerName, ps)
		}

	}
	return nil
}

// 处理chan
func HandleChan(psChan chan ProcessStatus) {
	for {
		ps := <-psChan

		prome.PromeRegister.Register(ps.PromeGauge)
		prome.PromeRegister.Register(ps.PromePIDGauge)
		if ps.Alive == 0 {
			// 启动通知
			log.Printf("process up by HandleChan : %#v", ps)
		}
		if ps.Alive == 1 {
			// 挂掉告警
			log.Printf("process down by HandleChan : %#v", ps)
		}
	}
}

// 获取用来设置promeSet的float64值
func GetGaugeValue(counter int64, pid string) float64 {
	gaugeValue, _ := strconv.ParseFloat(fmt.Sprintf("%d.%s", counter, pid), 64)
	return gaugeValue
}
