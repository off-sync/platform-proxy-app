package startproxy

func runProxy(c *Command, model *Model) {
	select {
	case <-model.Ctx.Done():
		return
	}
}
