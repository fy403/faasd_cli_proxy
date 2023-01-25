package api

type TestHandler struct{}

func (*TestHandler) Info(wait *WaitConn) {
	defer func() { wait.Done() }()
	wait.SetResult("", "test ok")
}
