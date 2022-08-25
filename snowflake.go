package idworker

import (
	"errors"
	"fmt"
	"hash/crc32"
	"math/rand"
	"net"
	"os"
	"sync"
	"time"

	"github.com/god-pig/gateway"
)

const (
	/**
	 * 时间起始标记点，作为基准，一般取系统的最近时间（一旦确定不能变动）
	 */
	twepoch = int64(1288834974657)
	/**
	* 机器标识位数
	 */
	workerIdBits     = int64(5)
	datacenterIdBits = int64(5)
	maxWorkerId      = int64(-1 ^ (-1 << workerIdBits))
	maxDatacenterId  = int64(-1 ^ (-1 << datacenterIdBits))
	/**
	* 毫秒内自增位
	 */
	sequenceBits      = int64(12)
	workerIdShift     = sequenceBits
	datacenterIdShift = sequenceBits + workerIdBits
	/**
	* 时间戳左移动位
	 */
	timestampLeftShift = sequenceBits + workerIdBits + datacenterIdBits
	sequenceMask       = int64(-1 ^ (-1 << sequenceBits))
)

type Snowflake struct {
	sync.Mutex
	lastTimestamp int64
	workerId      int64
	datacenterId  int64
	sequence      int64
}

func hashCode(str string) int {
	v := int(crc32.ChecksumIEEE([]byte(str)))
	if v >= 0 {
		return v
	}
	if -v >= 0 {
		return -v
	}
	// v == MinInt
	return 0
}

func getMaxWorkerId(datacenterId, maxWorkerId int64) int64 {
	mpid := fmt.Sprintf("%d%d", datacenterId, os.Getpid())

	/*
	 * MAC + PID 的 hashcode 获取16个低位
	 */
	code := int64(hashCode(mpid))
	return code & 0xffff % (maxWorkerId + 1)
}

func getInterfaceName(inetAddress string) (string, error) {
	if len(inetAddress) == 0 {
		g, err := gateway.GetGateway(gateway.IPv4)
		if err != nil {
			return "", err
		}
		return g.Interface, nil
	}
	interfaces, err := net.Interfaces()
	if err != nil {
		return "", err
	}
	for _, val := range interfaces {
		addrs, err := val.Addrs()
		if err != nil {
			continue
		}
		for _, addr := range addrs {
			if addr.String() == inetAddress {
				return val.Name, nil
			}
		}
	}
	return "", errors.New("gateway no found")
}

func getDatacenterId(maxDatacenterId int64, inetAddress string) int64 {
	id := int64(0)

	name, err := getInterfaceName(inetAddress)
	if err != nil {
		return 1
	}

	inet, err := net.InterfaceByName(name)
	if err != nil {
		return 0
	}

	mac := inet.HardwareAddr
	mac = mac[len(mac)-2:]

	id = ((0x000000FF & int64(mac[0])) | (0x0000FF00 & (int64(mac[1]) << 8))) >> 6
	id = id % (maxDatacenterId + 1)
	return id
}

func NewSnowflake(inetAddress string) *Snowflake {
	datacenterId := getDatacenterId(maxDatacenterId, inetAddress)
	workerId := getMaxWorkerId(datacenterId, maxDatacenterId)
	return &Snowflake{
		datacenterId:  datacenterId,
		workerId:      workerId,
		lastTimestamp: -1,
		sequence:      0,
	}
}

func (s *Snowflake) NextId() (string, error) {
	s.Lock()
	defer s.Unlock()
	timestamp := timeGen()
	//闰秒
	if timestamp < s.lastTimestamp {
		offset := s.lastTimestamp - timestamp
		if offset <= 5 {
			timestamp = timeGen()
			if timestamp < s.lastTimestamp {
				return "", fmt.Errorf("clock moved backwards. refusing to generate id for %d milliseconds", offset)
			}
		} else {
			return "", fmt.Errorf("clock moved backwards. refusing to generate id for %d milliseconds", offset)
		}
	}
	if s.lastTimestamp == timestamp {
		// 相同毫秒内，序列号自增
		s.sequence = (s.sequence + 1) & sequenceMask
		if s.sequence == 0 {
			// 同一毫秒的序列数已经达到最大
			timestamp = tilNextMillis(s.lastTimestamp)
		}
	} else {
		// 不同毫秒内，序列号置为 1 - 3 随机数
		s.sequence = int64(rand.Intn(2) + 1)
	}

	s.lastTimestamp = timestamp

	// 时间戳部分 | 数据中心部分 | 机器标识部分 | 序列号部分
	id := (timestamp-twepoch)<<timestampLeftShift |
		s.datacenterId<<datacenterIdShift |
		s.workerId<<workerIdShift |
		s.sequence
	return fmt.Sprintf("%d", id), nil
}

func tilNextMillis(lastTimestamp int64) int64 {
	timestamp := timeGen()
	for timestamp <= lastTimestamp {
		timestamp = timeGen()
	}
	return timestamp
}

func timeGen() int64 {
	return time.Now().UnixMilli()
}
