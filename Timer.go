package main

import (
	"container/list"
	"errors"
	"sync"
	"time"
)

type TimelineBuilder struct {
	TickSpeed time.Duration
	events    []*TimedEvent
}

func NewTimelineBuilder(tickSpeed time.Duration) *TimelineBuilder {
	return &TimelineBuilder{
		TickSpeed: tickSpeed,
	}
}

func (builder *TimelineBuilder) AddEvent(after time.Duration, callback func(*TickEvent)) {
	builder.events = append(builder.events, &TimedEvent{
		After:    after,
		Callback: callback,
	})
}

func (builder *TimelineBuilder) Build() *Timeline {
	return newTimeline(builder.TickSpeed, builder.events)
}

type TickEvent struct {
	Sender      *Timeline
	RunningTime time.Duration
	TimeLeft    time.Duration
}

type FinishedEvent struct {
	Sender   Timeline
	Duration time.Duration
	Time     time.Time
}

type TimedEvent struct {
	After    time.Duration
	Callback func(*TickEvent)
}

type Timeline struct {
	Tick     chan *TickEvent
	Finished chan *FinishedEvent

	startedAt time.Time
	tickSpeed time.Duration
	duration  time.Duration

	events         []*TimedEvent
	mutex          sync.Mutex
	eventsLen      int
	nextEventIndex int

	stop       chan bool
	isRunning  bool
	isFinished bool
}

func ToSlice(list *list.List) []interface{} {
	size := list.Len()
	result := make([]interface{}, size)

	var i int
	for e := list.Front(); e != nil; e = e.Next() {
		result[i] = e.Value
		i++
	}

	return result
}

func newTimeline(tickSpeed time.Duration, events []*TimedEvent) *Timeline {
	sortedEvents := list.New()
	sortedEvents.PushFront(events[0])

	for i := 1; i < len(events); i++ {
		toAdd := events[i]

		mark := sortedEvents.Back()
		for e := sortedEvents.Back(); e != nil && e.Value.(*TimedEvent).After > toAdd.After; e = e.Prev() {
			mark = e
		}

		if mark.Value.(*TimedEvent).After > toAdd.After {
			sortedEvents.PushFront(toAdd)
		} else {
			sortedEvents.PushBack(toAdd)
		}
	}

	slice := make([]*TimedEvent, sortedEvents.Len())
	untyped := ToSlice(sortedEvents)
	for i := 0; i < len(untyped); i++ {
		slice[i] = untyped[i].(*TimedEvent)
	}

	return &Timeline{
		Tick:     make(chan *TickEvent),
		Finished: make(chan *FinishedEvent),
		duration: sortedEvents.Back().Value.(*TimedEvent).After,

		events:         slice,
		eventsLen:      len(slice),
		nextEventIndex: 0,
		tickSpeed:      tickSpeed,
		stop:           make(chan bool),
		isRunning:      false,
		isFinished:     false,
	}
}

func (t *Timeline) Start() {
	if t.isRunning {
		errors.New("schedular is already running")
	}

	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.isRunning = true
	t.isFinished = false
	t.startedAt = time.Now()
	t.nextEventIndex = 0

	go t.run()
}

func (t *Timeline) run() {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	for !t.isFinished {
		select {
		case <-t.stop:
			return
		case <-time.After(t.tickSpeed):
			t.tick()
		}
	}

	if t.isFinished {
		t.Finished <- &FinishedEvent{}
	}
}

func (t *Timeline) Reset() {
	t.mutex.Lock()
	defer t.mutex.Unlock()

	t.Close()
	t.Start()
}

func (t *Timeline) tick() {
	now := time.Now()
	duration := now.Sub(t.startedAt)
	event := &TickEvent{
		Sender:      t,
		RunningTime: duration,
		TimeLeft:    t.duration - duration,
	}
	t.Tick <- event

	nextEvent := t.events[t.nextEventIndex]
	if nextEvent.After < duration {
		nextEvent.Callback(event)
		t.moveToNextEvent()
	}
}

func (t *Timeline) moveToNextEvent() {
	t.nextEventIndex++
	t.isFinished = t.nextEventIndex == t.eventsLen

	println("is finished: ", t.isFinished)
}

func (t *Timeline) Close() {
	t.stop <- true
	t.mutex.Lock()
	t.mutex.Unlock()

	t.isRunning = false
}
