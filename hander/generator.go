package hander

import (
	"fmt"
	"sync"
	"time"
)

const (
	workerIDBits   = uint64(10)                              // 节点ID的位数
	maxWorkerID    = int64(-1) ^ (int64(-1) << workerIDBits) // 节点ID的最大值
	sequenceBits   = uint64(12)                              // 序列号的位数
	workerIDShift  = sequenceBits                            // 节点ID左移位数
	timestampShift = sequenceBits + workerIDBits             // 时间戳左移位数
	sequenceMask   = int64(-1) ^ (int64(-1) << sequenceBits) // 序列号掩码
)

type Snowflake struct {
	mutex     sync.Mutex // 锁
	timestamp int64      // 上次生成ID的时间戳
	workerID  int64      // 节点ID
	sequence  int64      // 序列号
}

func NewSnowflake(workerID int64) (*Snowflake, error) {
	if workerID < 0 || workerID > maxWorkerID {
		return nil, fmt.Errorf("worker ID must be between 0 and %d", maxWorkerID)
	}
	return &Snowflake{
		timestamp: 0,
		workerID:  workerID,
		sequence:  0,
	}, nil
}

func (s *Snowflake) Generate() int64 {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	now := time.Now().UnixMilli()
	if s.timestamp == now {
		s.sequence = (s.sequence + 1) & sequenceMask
		if s.sequence == 0 {
			for now <= s.timestamp {
				now = time.Now().UnixMilli()
			}
		}
	} else {
		s.sequence = 0
	}

	s.timestamp = now

	id := ((now - 1288834974657) << timestampShift) |
		(s.workerID << workerIDShift) |
		(s.sequence)

	return id
}
