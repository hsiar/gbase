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

func (c *StrategyContext) Set(strategy Strategy) {
	c.strategy = strategy
}

func (c *StrategyContext) Exec(result any, args ...any) error {
	return c.strategy.Process(result, args...)
}

/*
usage:
  strategyContext := &StrategyContext{}
  strategyContext.Set(&StrategyA{})
  strategyContext.Exec(&result, arg1, arg2)
*/
