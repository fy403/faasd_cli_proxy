package api

import (
	"bytes"
	"fmt"
	"os/exec"
)

type FassCliHandler struct{}

func (*FassCliHandler) New(wait *WaitConn, req struct {
	lang   string `json:"lang" validate:"required"`
	name   string `json:"name" validate:"required"`
	prefix string `json:"prefix" validate:"required"`
}) {
	defer func() { wait.Done() }()
	if err := validate.Struct(&req); err != nil {
		wait.SetResult(err.Error(), "")
		return
	}
	// 移除同名文件，目录，避免数据干扰
	{
		rmCmd := exec.Command("rm", "-rf", req.name+"*")
		rmCmd.Start()
	}
	cmd := exec.Command("faas-cli", "new", "--lang", req.lang, req.name, "--prefix", req.prefix)
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		wait.SetResult(fmt.Sprintf("cmd exec failed: %s", err.Error()), "")
		return
	}
	wait.SetResult("", out.String())
}

func (*FassCliHandler) Write(wait *WaitConn, req struct {
	name    string `json:"name" validate:"required"`
	content []byte `json:"content"`
}) {
	defer func() { wait.Done() }()
	if err := validate.Struct(&req); err != nil {
		wait.SetResult(err.Error(), "")
		return
	}
}

func (*FassCliHandler) Build(wait *WaitConn, req struct {
	name string `json:"name" validate:"required"`
}) {
	defer func() { wait.Done() }()
	if err := validate.Struct(&req); err != nil {
		wait.SetResult(err.Error(), "")
		return
	}

	cmd := exec.Command("faas-cli", "build", "-f", req.name+".yml")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		wait.SetResult(fmt.Sprintf("cmd exec failed: %s", err.Error()), "")
		return
	}
	wait.SetResult("", out.String())
}

func (*FassCliHandler) Push(wait *WaitConn, req struct {
	name string `json:"name" validate:"required"`
}) {
	defer func() { wait.Done() }()
	if err := validate.Struct(&req); err != nil {
		wait.SetResult(err.Error(), "")
		return
	}

	cmd := exec.Command("faas-cli", "push", "-f", req.name+".yml")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		wait.SetResult(fmt.Sprintf("cmd exec failed: %s", err.Error()), "")
		return
	}
	wait.SetResult("", out.String())
}

func (*FassCliHandler) Deploy(wait *WaitConn, req struct {
	name string `json:"name" validate:"required"`
}) {
	defer func() { wait.Done() }()
	if err := validate.Struct(&req); err != nil {
		wait.SetResult(err.Error(), "")
		return
	}

	cmd := exec.Command("faas-cli", "deploy", "-f", req.name+".yml")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		wait.SetResult(fmt.Sprintf("cmd exec failed: %s", err.Error()), "")
		return
	}
	wait.SetResult("", out.String())
}

func (*FassCliHandler) GetInvokeInfo(wait *WaitConn, req struct {
	name string `json:"name" validate:"required"`
}) {
	defer func() { wait.Done() }()
	if err := validate.Struct(&req); err != nil {
		wait.SetResult(err.Error(), "")
		return
	}

	cmd := exec.Command("faas-cli")
	var out bytes.Buffer
	cmd.Stdout = &out
	if err := cmd.Run(); err != nil {
		wait.SetResult(fmt.Sprintf("cmd exec failed: %s", err.Error()), "")
		return
	}
	wait.SetResult("", out.String())
}
