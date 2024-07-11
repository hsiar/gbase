package gbase

//通用策略模式

// 策略接口
type Strategy interface {

	//result must be ptr
	Process(result any, args ...any) error
}

// 策略上下文
type StrategyContext struct {
	strategy Strategy
}

func (c *StrategyContext) SetStrategy(strategy Strategy) {
	c.strategy = strategy
}

func (c *StrategyContext) ExecuteStrategy(result any, args ...any) error {
	return c.strategy.Process(result, args...)
}

/*
usage:
  strategyContext := &StrategyContext{}
  strategyContext.SetStrategy(&StrategyA{})
  strategyContext.ExecuteStrategy(&result, arg1, arg2)
*/
