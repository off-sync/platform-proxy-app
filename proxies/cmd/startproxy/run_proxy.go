package startproxy

func runProxy(c *Command, model *Model) {
	select {
	case <-model.Ctx.Done():
		c.logger.Info("context is done: returning")
		return
	}
}
