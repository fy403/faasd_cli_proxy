package api

import (
	"fmt"
	"os/exec"
	"strings"

	"proxy/utils"
)

type FassCliHandler struct{}

type Data struct {
	Output       string `json:"output"`
	Dependencies string `json:"dependencies"`
	Content      string `json:"content"`
}

var langs = []string{
	"go",
	"golang-middleware",
	"golang-http",
	"node",
	"python",
	"python3",
}

func (*FassCliHandler) SupportedLang(wait *WaitConn) {
	defer func() {
		wait.Done()
	}()
	wait.SetResult("", langs)
}

func (*FassCliHandler) New(wait *WaitConn, req struct {
	Lang   string `json:"lang" validate:"required"`
	Name   string `json:"name" validate:"required"`
	Prefix string `json:"prefix" validate:"required"`
}) {
	defer func() { wait.Done() }()
	if err := validate.Struct(&req); err != nil {
		wait.SetResult(err.Error(), "")
		return
	}
	// 移除同名文件，目录，避免数据干扰
	{
		rmCmd := exec.Command("rm", "-rf", "./"+req.Name+"*")
		rmCmd.Run()
	}

	out, err := utils.ExecCommand("faas-cli", "new", "--lang", req.Lang, "--prefix", req.Prefix, req.Name)
	if err != nil {
		wait.SetResult(fmt.Sprintf("cmd exec failed: %s", err.Error()), "")
		return
	}

	var filePathWithDep, filePathWithCtn string
	switch req.Lang {
	case "go", "golang-middleware", "golang-http":
		// 写入./req.Name/go.mod
		filePathWithDep = fmt.Sprintf("./%s/go.mod", req.Name)
		filePathWithCtn = fmt.Sprintf("./%s/handler.go", req.Name)
	case "python", "python3":
		// 写入./req.Name/requirements.txt
		filePathWithDep = fmt.Sprintf("./%s/requirements.txt", req.Name)
		filePathWithCtn = fmt.Sprintf("./%s/handler.py", req.Name)
	case "node":
		// 写入./req.Name/package.json
		filePathWithDep = fmt.Sprintf("./%s/package.json", req.Name)
		filePathWithCtn = fmt.Sprintf("./%s/handler.js", req.Name)
	default:
		wait.SetResult("the lang: %s not supported", req.Lang)
		return
	}
	dep, err := utils.ReadFile(filePathWithDep)
	if err != nil {
		wait.SetResult(fmt.Sprintf("the file: %s not be read, err: %s", filePathWithDep, err.Error()), "")
		return
	}
	cnt, err := utils.ReadFile(filePathWithCtn)
	if err != nil {
		wait.SetResult(fmt.Sprintf("the file: %s not be read, err: %s", filePathWithCtn, err.Error()), "")
		return
	}
	// 返回模板数据
	data := Data{
		Output:       out,
		Dependencies: dep,
		Content:      cnt,
	}
	wait.SetResult("", data)
}

func (*FassCliHandler) Write(wait *WaitConn, req struct {
	Name         string `json:"name" validate:"required"`
	Lang         string `json:"lang" validate:"required"`
	Content      []byte `json:"content"`
	Dependencies []byte `json:"dependencies"`
}) {
	defer func() { wait.Done() }()
	if err := validate.Struct(&req); err != nil {
		wait.SetResult(err.Error(), "")
		return
	}
	// 按语言分别添加依赖
	var filePathWithDep, filePathWithCtn string
	switch req.Lang {
	case "go", "golang-middleware":
		// 写入./req.Name/go.mod
		filePathWithDep = fmt.Sprintf("./%s/go.mod", req.Name)
		filePathWithCtn = fmt.Sprintf("./%s/handler.go", req.Name)
	case "python", "python3":
		// 写入./req.Name/requirements.txt
		filePathWithDep = fmt.Sprintf("./%s/requirements.txt", req.Name)
		filePathWithCtn = fmt.Sprintf("./%s/handler.py", req.Name)
	case "node":
		// 写入./req.Name/package.json
		filePathWithDep = fmt.Sprintf("./%s/package.json", req.Name)
		filePathWithCtn = fmt.Sprintf("./%s/handler.js", req.Name)
	default:
		wait.SetResult("the lang: %s not supported", req.Lang)
		return
	}
	// 依赖写入
	if err := utils.WriteFile(filePathWithDep, req.Dependencies); err != nil {
		wait.SetResult(fmt.Sprintf("Could not write content to the file: %s, err: %s", filePathWithDep, err.Error()), "")
		return
	}
	// 代码写入
	if err := utils.WriteFile(filePathWithCtn, req.Content); err != nil {
		wait.SetResult(fmt.Sprintf("Could not write content to the file: %s, err: %s", filePathWithCtn, err.Error()), "")
		return
	}
	wait.SetResult("", "")
}

// build & push & deploy
func (*FassCliHandler) Up(wait *WaitConn, req struct {
	Name string `json:"name" validate:"required"`
}) {
	defer func() { wait.Done() }()
	if err := validate.Struct(&req); err != nil {
		wait.SetResult(err.Error(), "")
		return
	}
	out, err := utils.ExecCommand("faas-cli", "up", "-f", req.Name+".yml")
	if err != nil {
		wait.SetResult(fmt.Sprintf("cmd exec failed: %s", err.Error()), "")
		return
	}
	wait.SetResult("", out)
}

func (*FassCliHandler) Build(wait *WaitConn, req struct {
	Name string `json:"name" validate:"required"`
}) {
	defer func() { wait.Done() }()
	if err := validate.Struct(&req); err != nil {
		wait.SetResult(err.Error(), "")
		return
	}

	out, err := utils.ExecCommand("faas-cli", "build", "-f", req.Name+".yml")
	if err != nil {
		wait.SetResult(fmt.Sprintf("cmd exec failed: %s", err.Error()), "")
		return
	}
	wait.SetResult("", out)
}

func (*FassCliHandler) Push(wait *WaitConn, req struct {
	Name string `json:"name" validate:"required"`
}) {
	defer func() { wait.Done() }()
	if err := validate.Struct(&req); err != nil {
		wait.SetResult(err.Error(), "")
		return
	}

	out, err := utils.ExecCommand("faas-cli", "push", "-f", req.Name+".yml")
	if err != nil {
		wait.SetResult(fmt.Sprintf("cmd exec failed: %s", err.Error()), "")
		return
	}
	wait.SetResult("", out)
}

func (*FassCliHandler) Deploy(wait *WaitConn, req struct {
	Name string `json:"name" validate:"required"`
}) {
	defer func() { wait.Done() }()
	if err := validate.Struct(&req); err != nil {
		wait.SetResult(err.Error(), "")
		return
	}

	out, err := utils.ExecCommand("faas-cli", "deploy", "-f", req.Name+".yml")
	if err != nil {
		wait.SetResult(fmt.Sprintf("cmd exec failed: %s", err.Error()), "")
		return
	}
	wait.SetResult("", out)
}

func (*FassCliHandler) GetAllInvokeInfo(wait *WaitConn) {
	defer func() { wait.Done() }()
	out, err := utils.ExecCommand("faas-cli", "list", "-q")
	if err != nil {
		wait.SetResult(fmt.Sprintf("cmd exec failed: %s", err.Error()), "")
		return
	}
	functions := strings.Split(strings.TrimSpace(out), "\n")
	var describes []string
	for _, function := range functions {
		out, err = utils.ExecCommand("faas-cli", "describe", function)
		if err != nil {
			wait.SetResult(fmt.Sprintf("cmd exec failed: %s", err.Error()), "")
			return
		}
		describes = append(describes, out)
	}
	wait.SetResult("", describes)
}

func (*FassCliHandler) GetInvokeInfo(wait *WaitConn, req struct {
	Name string `json:"name" validate:"required"`
}) {
	defer func() { wait.Done() }()
	if err := validate.Struct(&req); err != nil {
		wait.SetResult(err.Error(), "")
		return
	}

	out, err := utils.ExecCommand("faas-cli", "describe", req.Name)
	if err != nil {
		wait.SetResult(fmt.Sprintf("cmd exec failed: %s", err.Error()), "")
		return
	}
	wait.SetResult("", out)
}
