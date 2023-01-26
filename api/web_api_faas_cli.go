package api

import (
	"bytes"
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

	var filePathWithDep, filePathWithCtn string
	switch req.lang {
	case "go", "golang-middleware":
		// 写入./req.name/go.mod
		filePathWithDep = fmt.Sprintf("./%s/go.mod", req.name)
		filePathWithCtn = fmt.Sprintf("./%s/handler.go", req.name)
	case "python", "python3":
		// 写入./req.name/requirements.txt
		filePathWithDep = fmt.Sprintf("./%s/requirements.txt", req.name)
		filePathWithCtn = fmt.Sprintf("./%s/handler.py", req.name)
	case "node":
		// 写入./req.name/package.json
		filePathWithDep = fmt.Sprintf("./%s/package.json", req.name)
		filePathWithCtn = fmt.Sprintf("./%s/handler.js", req.name)
	default:
		wait.SetResult("the lang: %s not supported", req.lang)
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
		Output:       out.String(),
		Dependencies: dep,
		Content:      cnt,
	}
	wait.SetResult("", data)
}

func (*FassCliHandler) Write(wait *WaitConn, req struct {
	name         string `json:"name" validate:"required"`
	lang         string `json:"lang" validate:"required"`
	content      []byte `json:"content"`
	dependencies []byte `json:"dependencies"`
}) {
	defer func() { wait.Done() }()
	if err := validate.Struct(&req); err != nil {
		wait.SetResult(err.Error(), "")
		return
	}
	// 按语言分别添加依赖
	var filePathWithDep, filePathWithCtn string
	switch req.lang {
	case "go", "golang-middleware":
		// 写入./req.name/go.mod
		filePathWithDep = fmt.Sprintf("./%s/go.mod", req.name)
		filePathWithCtn = fmt.Sprintf("./%s/handler.go", req.name)
	case "python", "python3":
		// 写入./req.name/requirements.txt
		filePathWithDep = fmt.Sprintf("./%s/requirements.txt", req.name)
		filePathWithCtn = fmt.Sprintf("./%s/handler.py", req.name)
	case "node":
		// 写入./req.name/package.json
		filePathWithDep = fmt.Sprintf("./%s/package.json", req.name)
		filePathWithCtn = fmt.Sprintf("./%s/handler.js", req.name)
	default:
		wait.SetResult("the lang: %s not supported", req.lang)
		return
	}
	// 依赖写入
	if err := utils.WriteFile(filePathWithDep, req.dependencies); err != nil {
		wait.SetResult(fmt.Sprintf("Could not write content to the file: %s, err: %s", filePathWithDep, err.Error()), "")
		return
	}
	// 代码写入
	if err := utils.WriteFile(filePathWithCtn, req.content); err != nil {
		wait.SetResult(fmt.Sprintf("Could not write content to the file: %s, err: %s", filePathWithCtn, err.Error()), "")
		return
	}
	wait.SetResult("", "")
}

// build & push & deploy
func (*FassCliHandler) Up(wait *WaitConn, req struct {
	name string `json:"name" validate:"required"`
}) {
	defer func() { wait.Done() }()
	if err := validate.Struct(&req); err != nil {
		wait.SetResult(err.Error(), "")
		return
	}
	out, err := utils.ExecCommand("faas-cli", "up", "-f", req.name+".yml")
	if err != nil {
		wait.SetResult(fmt.Sprintf("cmd exec failed: %s", err.Error()), "")
		return
	}
	wait.SetResult("", out)
}

func (*FassCliHandler) Build(wait *WaitConn, req struct {
	name string `json:"name" validate:"required"`
}) {
	defer func() { wait.Done() }()
	if err := validate.Struct(&req); err != nil {
		wait.SetResult(err.Error(), "")
		return
	}

	out, err := utils.ExecCommand("faas-cli", "build", "-f", req.name+".yml")
	if err != nil {
		wait.SetResult(fmt.Sprintf("cmd exec failed: %s", err.Error()), "")
		return
	}
	wait.SetResult("", out)
}

func (*FassCliHandler) Push(wait *WaitConn, req struct {
	name string `json:"name" validate:"required"`
}) {
	defer func() { wait.Done() }()
	if err := validate.Struct(&req); err != nil {
		wait.SetResult(err.Error(), "")
		return
	}

	out, err := utils.ExecCommand("faas-cli", "push", "-f", req.name+".yml")
	if err != nil {
		wait.SetResult(fmt.Sprintf("cmd exec failed: %s", err.Error()), "")
		return
	}
	wait.SetResult("", out)
}

func (*FassCliHandler) Deploy(wait *WaitConn, req struct {
	name string `json:"name" validate:"required"`
}) {
	defer func() { wait.Done() }()
	if err := validate.Struct(&req); err != nil {
		wait.SetResult(err.Error(), "")
		return
	}

	out, err := utils.ExecCommand("faas-cli", "deploy", "-f", req.name+".yml")
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
	name string `json:"name" validate:"required"`
}) {
	defer func() { wait.Done() }()
	if err := validate.Struct(&req); err != nil {
		wait.SetResult(err.Error(), "")
		return
	}

	out, err := utils.ExecCommand("faas-cli", "describe", req.name)
	if err != nil {
		wait.SetResult(fmt.Sprintf("cmd exec failed: %s", err.Error()), "")
		return
	}
	wait.SetResult("", out)
}
