package gbase

//通用策略模式

// 策略接口
//type Strategy interface {
//
//	//result must be ptr
//	Process(result any, args ...any) error
//}
//
//
//
//// 策略上下文
//type StrategyContext struct {
//	strategy Strategy
//}
//
//func (c *StrategyContext) Set(strategy Strategy) {
//	c.strategy = strategy
//}
//
//func (c *StrategyContext) Exec(result any, args ...any) error {
//	return c.strategy.Process(result, args...)
//}

type IStrategy[T any] interface {
	// Process 方法返回 T 类型的值
	Process(args ...any) (T, error)
	// 获取类型，用于区分不同的策略
	GetType() int
	// 设置及获取数据，用于外界数据传入策略逻辑
	SetData(data ...any)
}

type BaseStrategy struct{}                  //空策略
func (s *BaseStrategy) GetType() int        { return 0 }
func (s *BaseStrategy) SetData(data ...any) {}
