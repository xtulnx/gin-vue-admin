package ebus

import "github.com/asaskevich/EventBus"

// EventBus 事件总线
//
// jason.liao 2023.12.08

var Bus EventBus.Bus = EventBus.New()

var disableAll bool = false

// InitSys 系统初始化，这里不要依赖任何组件。
func InitSys() {

}

// DisableAll 禁用所有事件
func DisableAll() {
	disableAll = true
}

// BlockTopic 屏蔽一些事件
func BlockTopic(topics ...string) {
	// todo
}

// 是否屏蔽某个事件
func kitIsBlock(topic string) bool {
	if disableAll {
		return true
	}
	return false
}

// Publish 发布事件, 事件名，参数（参数与订阅的函数参数一致，需要事先约定好）
func Publish(topic string, args ...interface{}) {
	if kitIsBlock(topic) {
		return
	}
	Bus.Publish(topic, args...)
}

// Subscribe 订阅事件（同步）
func Subscribe(topic string, fn interface{}) error {
	return Bus.Subscribe(topic, fn)
}

// SubscribeAsync 订阅事件（异步）
func SubscribeAsync(topic string, fn interface{}, transactional bool) error {
	return Bus.SubscribeAsync(topic, fn, transactional)
}

// SubscribeOnce 订阅一次事件（触发后就会删除订阅）(同步)
func SubscribeOnce(topic string, fn interface{}) error {
	return Bus.SubscribeOnce(topic, fn)
}

// SubscribeOnceAsync 订阅一次事件（触发后就会删除订阅）(异步)
func SubscribeOnceAsync(topic string, fn interface{}) error {
	return Bus.SubscribeOnceAsync(topic, fn)
}

// Unsubscribe 取消订阅（根据回调函数）
func Unsubscribe(topic string, handler interface{}) error {
	return Bus.Unsubscribe(topic, handler)
}

// HasCallback 检查指定 topic 是否有订阅者
func HasCallback(topic string) bool {
	return Bus.HasCallback(topic)
}

// WaitAsync 等待所有异步事件完成
func WaitAsync() {
	Bus.WaitAsync()
}
