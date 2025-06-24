package di

func (c *Container) Use(steps ...func(*Container) error) *Container {
	c.tasks = append(c.tasks, steps...)
	return c
}

func (c *Container) Run() error {
	for _, task := range c.tasks {
		if err := task(c); err != nil {
			return err
		}
	}
	return nil
}

func Provide(provider any) func(*Container) error {
	return func(c *Container) error {
		return c.Provide(provider)
	}
}

func Invoke(fn any) func(*Container) error {
	return func(c *Container) error {
		return c.Invoke(fn)
	}
}
