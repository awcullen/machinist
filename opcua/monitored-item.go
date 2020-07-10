package opcua

import (
	"time"

	"sync"

	ua "github.com/awcullen/opcua"
	deque "github.com/gammazero/deque"
)

const (
	maxQueueSize        = 1024
	minSamplingInterval = 100.0
	maxSamplingInterval = 60 * 1000.0
)

// MonitoredItem specifies the node that is monitored for data changes or events.
type MonitoredItem struct {
	sync.RWMutex
	metricID       string
	monitoringMode ua.MonitoringMode
	queueSize      uint32
	currentValue   *ua.DataValue
	queue          deque.Deque
	ts             time.Time
	ti             time.Duration
}

// NewMonitoredItem constructs a new MonitoredItem.
func NewMonitoredItem(metricID string, monitoringMode ua.MonitoringMode, samplingInterval float64, queueSize uint32, ts time.Time) *MonitoredItem {
	if samplingInterval > maxSamplingInterval {
		samplingInterval = maxSamplingInterval
	}
	if samplingInterval < minSamplingInterval {
		samplingInterval = minSamplingInterval
	}
	if queueSize > maxQueueSize {
		queueSize = maxQueueSize
	}
	if queueSize < 1 {
		queueSize = 1
	}
	mi := &MonitoredItem{
		metricID:       metricID,
		monitoringMode: monitoringMode,
		queueSize:      queueSize,
		ti:             time.Duration(samplingInterval) * time.Millisecond,
		queue:          deque.Deque{},
		currentValue:   ua.NewDataValueVariant(&ua.NilVariant, ua.BadWaitingForInitialData, time.Time{}, 0, time.Time{}, 0),
		ts:             ts,
	}
	return mi
}

// Enqueue stores the current value of the item.
func (mi *MonitoredItem) Enqueue(value *ua.DataValue) {
	mi.Lock()
	for mi.queue.Len() >= int(mi.queueSize) {
		mi.queue.PopFront() // discard oldest
	}
	mi.queue.PushBack(value)
	mi.Unlock()
}

// Notifications returns the data values from the queue, up to given time and max count.
func (mi *MonitoredItem) Notifications(tn time.Time, max int) (notifications []*ua.DataValue, more bool) {
	mi.Lock()
	defer mi.Unlock()
	notifications = make([]*ua.DataValue, 0, mi.queueSize)
	if mi.monitoringMode != ua.MonitoringModeReporting {
		return notifications, false
	}
	// if in sampling interval mode, append the last value of each sampling interval
	if mi.ti > 0 {
		if mi.ts.IsZero() {

		}
		// for each interval (a tumbling window that ends at ts (the sample time))
		for ; !mi.ts.After(tn) && len(notifications) < max; mi.ts = mi.ts.Add(mi.ti) {
			// for each value in queue
			for mi.queue.Len() > 0 {
				// peek
				peek := mi.queue.Front().(*ua.DataValue)
				// if peeked value timestamp is before ts (the sample time)
				if !peek.ServerTimestamp().After(mi.ts) {
					// overwrite current with the peeked value
					mi.currentValue = peek
					// pop it from the queue
					mi.queue.PopFront()
				} else {
					// peeked value timestamp is later
					break
				}
			}
			// update current's timestamp, and append it to the notifications to return
			cv := mi.currentValue
			mi.currentValue = ua.NewDataValueVariant(cv.InnerVariant(), cv.StatusCode(), cv.SourceTimestamp(), 0, mi.ts, 0)
			notifications = append(notifications, mi.currentValue)
		}
	} else {
		// for each value in queue
		for mi.queue.Len() > 0 && len(notifications) < max {
			peek := mi.queue.Front().(*ua.DataValue)
			// if peeked value timestamp is before tn (the sample time)
			if !peek.ServerTimestamp().After(tn) {
				// append it to the notifications to return
				notifications = append(notifications, peek)
				mi.currentValue = peek
				// pop it from the queue
				mi.queue.PopFront()
			} else {
				// peeked value timestamp is later
				break
			}
		}
	}

	more = mi.queue.Len() > 0
	return notifications, more
}
