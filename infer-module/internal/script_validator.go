package internal

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/MingkaiLee/kasos/infer-module/client"
	"github.com/MingkaiLee/kasos/infer-module/config"
	"github.com/MingkaiLee/kasos/infer-module/util"
	"github.com/google/uuid"
	jsoniter "github.com/json-iterator/go"
)

type ScriptValidateResult struct {
	ModelName string // 模型名
	Err       error  // 错误信息
}
type ScriptValidator struct {
	modelName   string
	trainScript *string
	inferScript *string
	result      chan *ScriptValidateResult
}

type ScriptValidateRequest struct {
	ModelName   string  `json:"model_name"`
	TrainScript *string `json:"train_script"`
	InferScript *string `json:"infer_script"`
}

func ScriptValidate(ctx context.Context, content []byte) (err error) {
	var req ScriptValidateRequest
	err = jsoniter.Unmarshal(content, &req)
	if err != nil {
		util.LogErrorf("internal.ScriptValidate: unmarshal error: %v", err)
		return
	}
	if req.ModelName == "" || req.TrainScript == nil || req.InferScript == nil {
		err = fmt.Errorf("invalid request")
		util.LogErrorf("internal.ScriptValidate: invalid request: %v", req)
		return
	}
	validator := NewScriptValidator(req.ModelName, req.TrainScript, req.InferScript)
	validator.Run()
	go func() {
		result := <-validator.GetResult()
		util.LogInfof("service.ScriptValidate result: %+v", result)
		var req client.ReportModelValidRequest
		req.ModelName = result.ModelName
		if result.Err != nil {
			req.Ok = false
			req.ErrorInfo = result.Err.Error()
		} else {
			req.Ok = true
		}
		resp, e := client.CallReportModelValid(ctx, &req)
		if e != nil {
			util.LogErrorf("service.ScriptValidate report request: %+v,resp: %+v, err: %v", req, *resp, e)
		}
		util.LogInfof("service.ScriptValidate report request: %+v,resp: %+v, err: %v", req, *resp, e)
	}()
	return
}

func NewScriptValidator(modelName string, trainScript, inferScript *string) *ScriptValidator {
	return &ScriptValidator{
		modelName:   modelName,
		trainScript: trainScript,
		inferScript: inferScript,
		result:      make(chan *ScriptValidateResult, 1),
	}
}

func (s *ScriptValidator) Run() {
	go func() {
		r := new(ScriptValidateResult)
		r.ModelName = s.modelName
		defer func() {
			s.result <- r
		}()
		// 临时测试脚本的目录
		tmpDir := fmt.Sprintf("%s/tmp", config.ScriptDirectory)
		// 创建临时训练脚本
		trainScriptFile, err := os.CreateTemp(tmpDir, "*.py")
		if err != nil {
			r.Err = err
			util.LogErrorf("internal.ScriptValidator.Run: create temp file error: %v", err)
			return
		}
		defer func() {
			// 关闭并删除临时文件, 遇到错误时忽略
			trainScriptFile.Close()
			os.Remove(trainScriptFile.Name())
		}()
		// 写入临时训练脚本
		_, err = trainScriptFile.WriteString(*s.trainScript)
		if err != nil {
			r.Err = err
			util.LogErrorf("internal.ScriptValidator.Run: write temp file error: %v", err)
			return
		}
		// 创建临时推理脚本
		inferScriptFile, err := os.CreateTemp(tmpDir, "*.py")
		if err != nil {
			r.Err = err
			util.LogErrorf("internal.ScriptValidator.Run: create temp file error: %v", err)
			return
		}
		defer func() {
			// 关闭并删除临时文件, 遇到错误时忽略
			inferScriptFile.Close()
			os.Remove(inferScriptFile.Name())
		}()
		// 写入临时推理脚本
		_, err = inferScriptFile.WriteString(*s.inferScript)
		if err != nil {
			r.Err = err
			util.LogErrorf("internal.ScriptValidator.Run: write temp file error: %v", err)
			return
		}
		// 首次训练测试
		// 创建一个模型文件的路径
		modelPath := fmt.Sprintf("%s/tmp/%s", config.ModelDirectory, uuid.NewString())
		// 结束时删除临时模型文件
		defer func() {
			os.Remove(modelPath)
		}()
		firstTrainCmd := exec.Command("python3", trainScriptFile.Name(), "--new", "-d", config.ValidateDataPath, "-m", modelPath)
		err = firstTrainCmd.Run()
		if err != nil {
			r.Err = err
			output, _ := firstTrainCmd.CombinedOutput()
			util.LogErrorf("internal.ScriptValidator.Run: command output: %s", string(output))
			util.LogErrorf("internal.ScriptValidator.Run: first train error: %v", err)
			return
		}
		// 迭代训练测试
		furtherTrainCmd := exec.Command("python3", trainScriptFile.Name(), "-d", config.ValidateDataPath, "-m", modelPath)
		err = furtherTrainCmd.Run()
		if err != nil {
			r.Err = err
			output, _ := furtherTrainCmd.CombinedOutput()
			util.LogErrorf("internal.ScriptValidator.Run: command output: %s", string(output))
			util.LogErrorf("internal.ScriptValidator.Run: further train error: %v", err)
			return
		}
		// 推理测试
		inferCmd := exec.Command("python3", inferScriptFile.Name(), "-t", time.DateTime, "-v", "6.6666", "-m", modelPath)
		err = inferCmd.Run()
		if err != nil {
			r.Err = err
			output, _ := inferCmd.CombinedOutput()
			util.LogErrorf("internal.ScriptValidator.Run: command output: %s", string(output))
			util.LogErrorf("internal.ScriptValidator.Run: infer error: %v", err)
			return
		}
		// 全部测试完成后可正式将脚本写入脚本目录中, 直接通过复制功能完成
		trainFormalPath := fmt.Sprintf("%s/train/%s.py", config.ScriptDirectory, s.modelName)
		inferFormalPath := fmt.Sprintf("%s/infer/%s.py", config.ScriptDirectory, s.modelName)
		trainFormalFile, err := os.Create(trainFormalPath)
		if err != nil {
			r.Err = err
			util.LogErrorf("internal.ScriptValidator.Run: create formal file error: %v", err)
			return
		}
		defer trainFormalFile.Close()
		inferFormalFile, err := os.Create(inferFormalPath)
		if err != nil {
			r.Err = err
			util.LogErrorf("internal.ScriptValidator.Run: create formal file error: %v", err)
			return
		}
		defer inferFormalFile.Close()
		_, err = trainFormalFile.WriteString(*s.trainScript)
		if err != nil {
			r.Err = err
			util.LogErrorf("internal.ScriptValidator.Run: write train script error: %v", err)
			return
		}
		_, err = inferFormalFile.WriteString(*s.inferScript)
		if err != nil {
			r.Err = err
			util.LogErrorf("internal.ScriptValidator.Run: write infer script error: %v", err)
			return
		}
	}()
}

// 验证结果的channel
func (s *ScriptValidator) GetResult() <-chan *ScriptValidateResult {
	return s.result
}
