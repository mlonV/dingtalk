package dnsquery

import (
	"context"
	"net"
	"sync"
	"time"

	"github.com/miekg/dns"
	"github.com/mlonV/dingtalk/config"
	"github.com/mlonV/dingtalk/prome"
	"github.com/prometheus/client_golang/prometheus"
)

type Metrics struct {
	DnsGauge prometheus.Gauge

	queryType uint16 // A or CNAME
	server    string
	domain    string
	cname     string
	arecord   string

	lastCNAME string
	lastA     string

	// 计数器，避免快速恢复
	changeCounter int

	// 使用channel来取消
	unRegisterCh chan struct{}

	// 锁来取消  延迟取消注册
	mu sync.Mutex
	// 标志位，避免重复注册
	isRegistering bool
}

var (
	serverCache = make(map[string]map[string]*Metrics)
	// { server : { domain : cname }}
)

func init() {
	for _, server := range config.Conf.DnsQuery.Servers {
		serverCache[server] = make(map[string]*Metrics)
		// 查询A记录的
		for _, domain := range config.Conf.DnsQuery.ARecords {
			m := NewMetrics(server, domain, "", "", dns.TypeA)
			m.queryDNSForMetrics()
			m.NewGauge()
			serverCache[server][domain] = m
		}
		// 查询CNAME的
		for _, domain := range config.Conf.DnsQuery.CnameRecords {
			m := NewMetrics(server, domain, "", "", dns.TypeCNAME)
			m.queryDNSForMetrics()
			m.NewGauge()
			serverCache[server][domain] = m
		}
	}

}

func NewMetrics(server, domain, cname string, arecord string, t uint16) *Metrics {
	return &Metrics{
		queryType:     t,
		server:        server,
		domain:        domain,
		arecord:       arecord,
		cname:         cname,
		changeCounter: 0,
		unRegisterCh:  nil,
	}
}

func (m *Metrics) NewGauge() {

	// m.DnsGauge = prometheus.NewGaugeVec(
	// 	prometheus.GaugeOpts{
	// 		Namespace: "karawangroup",
	// 		Subsystem: "dns",
	// 		Name:      "dnsquery",
	// 		Help:      "DNS change detection",
	// 	},
	// 	[]string{"server", "domain", "Arecord", "CNAME"},
	// )

	m.DnsGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: "karawangroup",
		Subsystem: "dns",
		Name:      "dnsquery",
		Help:      "Indicates if there is a change in CNAME record (1 for change, 0 for no change)",
		ConstLabels: map[string]string{
			"server":  m.server,
			"domain":  m.domain + config.Conf.DnsQuery.Domain,
			"arecord": m.arecord,
			"CNAME":   m.cname,
		},
	})
}

func (m *Metrics) Unregister() bool {
	return prome.PromeRegister.Unregister(m.DnsGauge)
}
func (m *Metrics) Register() {
	prome.PromeRegister.Register(m.DnsGauge)
}

func (m *Metrics) queryDNSForMetrics() {
	msg := new(dns.Msg)
	msg.SetQuestion(dns.Fqdn(m.domain+config.Conf.DnsQuery.Domain), m.queryType)

	client := new(dns.Client)
	client.Timeout = time.Duration(config.Conf.DnsQuery.QueryTimeOut) * time.Second
	in, _, err := client.Exchange(msg, net.JoinHostPort(m.server, "53"))
	if err != nil {
		// fmt.Println("DNS query failed for %s on %s in.Rcode : %s   err: %s ", domain, server, in.Rcode, err.Error())
		config.Log.Error("DNS query failed for %s on %s   err:  ", m.domain, m.server, err)
		return
	}

	if in.Rcode != dns.RcodeSuccess {
		config.Log.Error("DNS query failed for %s on %s   err:  ", m.domain, m.server, err)
		return
	}

	for _, answer := range in.Answer {
		if m.queryType == dns.TypeA {
			if aRecord, ok := answer.(*dns.A); ok {
				m.arecord = aRecord.A.String()
				config.Log.Info("domainserver: %s domain : %s   A record: %s\n", m.server, m.domain, aRecord.A.String())
			}
		} else if m.queryType == dns.TypeCNAME {
			if cnameRecord, ok := answer.(*dns.CNAME); ok {
				m.cname = cnameRecord.Target
				config.Log.Info("domainserver: %s domain : %s  CName record: %#v\n", m.server, m.domain, cnameRecord.Target)

			}
		}
	}
}

// 使用channel来取消注册
func (m *Metrics) checkWithChannel() {
	m.queryDNSForMetrics()
	m.Register()
	if m.lastCNAME == "" && m.lastA == "" {
		m.lastCNAME = m.cname
		m.lastA = m.arecord
	}

	if m.lastCNAME != m.cname || m.lastA != m.arecord {
		m.changeCounter++
		if m.changeCounter >= 5 {

			// 延迟取消旧指标（异步安全）
			oldGauge := m.DnsGauge
			if m.unRegisterCh == nil {
				m.unRegisterCh = make(chan struct{})
				go func(m *Metrics) {
					time.Sleep(time.Duration(config.Conf.DnsQuery.DelayUnregister) * time.Second)

					res := prome.PromeRegister.Unregister(oldGauge)
					config.Log.Info("unregister result: %v  Server: %s  domain : %s  CName : %s A: %s\n", res, m.server, m.domain, m.cname, m.arecord, m.DnsGauge)
					// 重新注册一次（解析的IP变化了

					m.unRegisterCh = nil
				}(m)
			}
			m.lastCNAME = m.cname
			m.lastA = m.arecord
			m.NewGauge()
			m.Register()
			m.changeCounter = 0

		}
		m.DnsGauge.Set(float64(m.changeCounter))

	} else {
		if m.changeCounter > 0 {
			m.changeCounter--
		}
		if m.changeCounter == 0 {
			m.DnsGauge.Set(float64(m.changeCounter))
		}
	}

}

// 使用锁来取消注册
func (m *Metrics) check() {
	m.queryDNSForMetrics()
	m.Register()
	if m.lastCNAME == "" && m.lastA == "" {
		m.lastCNAME = m.cname
		m.lastA = m.arecord
	}

	if m.lastCNAME != m.cname || m.lastA != m.arecord {
		m.changeCounter++
		if m.changeCounter >= 5 {

			// 延迟取消旧指标（异步安全）
			if !m.isRegistering {
				m.isRegistering = true
				go func(oldGauge prometheus.Gauge) {
					time.Sleep(time.Duration(config.Conf.DnsQuery.DelayUnregister) * time.Second)
					m.mu.Lock()
					defer m.mu.Unlock()
					res := prome.PromeRegister.Unregister(oldGauge)
					config.Log.Info("unregister result: %v  Server: %s  domain : %s  CName : %s A: %s\n", res, m.server, m.domain, m.cname, m.arecord, m.DnsGauge)
					m.isRegistering = false
				}(m.DnsGauge)
			}

			// 重新注册一次（解析的IP变化了
			m.lastCNAME = m.cname
			m.lastA = m.arecord
			m.NewGauge()
			m.Register()
			m.changeCounter = 0
		}
		m.DnsGauge.Set(float64(m.changeCounter))

	} else {
		if m.changeCounter > 0 {
			m.changeCounter--
		}
		if m.changeCounter == 0 {
			m.DnsGauge.Set(float64(m.changeCounter))
		}
	}

}

func worker(ctx context.Context) {
	config.Log.Info("Starting DNS Query Worker, Interval: %d s", config.Conf.DnsQuery.Interval)
	interval := time.Duration(config.Conf.DnsQuery.Interval) * time.Second
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			for _, domainMap := range serverCache {
				for _, metrics := range domainMap {
					go metrics.check()
				}
			}
		}
	}
}

func StartDNSQueryWorker(ctx context.Context) {
	go worker(ctx)
}
