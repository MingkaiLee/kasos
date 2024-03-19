package internal

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"strconv"

	"github.com/MingkaiLee/kasos/infer-module/client"
	"github.com/MingkaiLee/kasos/infer-module/config"
	"github.com/MingkaiLee/kasos/infer-module/util"
)

type Service struct {
	serviceName string // 服务名称
	modelName   string // 模型名称
	tags        string // 服务标签
	scriptPath  string // 推理脚本的路径
	modelPath   string // 模型的路径
	// 下方内容用于没有时序模型时的预测逻辑
	prediction int // 上一周期的预测值
	diff       int // 上一周期的实际值与预测值的差
}

func NewService(serviceName, modelName, tags string) *Service {
	return &Service{
		serviceName: serviceName,
		modelName:   modelName,
		tags:        tags,
		scriptPath:  fmt.Sprintf("%s/infer/%s.py", config.ScriptDirectory, modelName),
		modelPath:   fmt.Sprintf("%s/%s/%s", config.ModelDirectory, serviceName, modelName),
		prediction:  -1, // 初始值为-1表示历史上从未作出预测
	}
}

func (s *Service) Run(ctx context.Context) (err error) {
	// 获取最新的数据
	data, err := client.RealTimeData(ctx, s.serviceName)
	if err != nil {
		util.LogErrorf("service.Service.Run error: %v", err)
		return
	}
	// 推理
	prediction, err := s.infer(data)
	if err != nil {
		util.LogErrorf("service.Service.Run error: %v", err)
		return
	}
	// 上报至HPA-EXECUTOR
	req := client.ReportQPSRequest{
		ServiceName: s.serviceName,
		QPS:         prediction,
	}
	resp, err := client.CallReportQPS(ctx, &req)
	if err != nil {
		util.LogErrorf("service.Service.Run error: %v", err)
		return
	}
	if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("report qps error, status: %s, code: %d", resp.Status, resp.StatusCode)
		util.LogErrorf("service.Service.Run error: %v", err)
	}
	return
}

// 推理
func (s *Service) infer(data client.SerialDataPoint) (prediction int, err error) {
	// 检查模型是否已存在
	modelFileStat, err := os.Stat(s.modelPath)
	if os.IsNotExist(err) {
		// 模型不存在, 则使用预测逻辑
		// 因浮点数的截断问题, 保留1
		p := int(data.Value) + 1
		if s.prediction < 0 {
			// 历史从未作出预测时, 预测值同当前实际值
			prediction = p
			util.LogInfof("service.Service.infer, first prediction: %d", prediction)
			return
		}
		// 上一周期qps预测值与实际值的差异作diff
		diff := p - s.prediction
		// 对下一周期qps预测值作修正
		if diff > 0 {
			if s.diff > 0 {
				// 连续两个周期, 预测值低于实际值, 则累加两个周期的diff
				prediction = p + diff + s.diff
				s.diff = diff
				util.LogInfof("service.Service.infer, prediction: %d", prediction)
				return
			} else {
				// 上个周期的预测值高于实际值, 则只加本周期的diff
				prediction = p + diff
				s.diff = diff
				util.LogInfof("service.Service.infer, prediction: %d", prediction)
				return
			}
		} else {
			if s.diff > 0 {
				// 上一周期的预测值低于实际值, 则只加本周期的diff
				prediction = p + diff
				s.diff = diff
				util.LogInfof("service.Service.infer, prediction: %d", prediction)
				return
			} else {
				// 连续两个周期, 预测值高于实际值, 则累加两个周期的diff
				prediction = p + diff + s.diff
				s.diff = diff
				util.LogInfof("service.Service.infer, prediction: %d", prediction)
				return
			}
		}
	} else {
		// 模型文件若为目录, 则遇到了严重错误
		if modelFileStat.IsDir() {
			util.LogErrorf("internal.Service.infer, model is a directory, expected a file")
			return
		}
		// 模型已存在, 执行推理
		cmd := exec.Command("python3", s.scriptPath, "-t", data.Timestamp, "-v", strconv.FormatFloat(data.Value, 'f', 4, 64), "-m", s.modelPath)
		var out bytes.Buffer
		cmd.Stdout = &out
		err = cmd.Run()
		if err != nil {
			util.LogErrorf("internal.Service.infer, infer failed, error: %v", err)
			return
		}
		// 获取预测值
		result := out.String()
		// 将预测值转为float64类型
		p, e := strconv.ParseFloat(result, 64)
		if e != nil {
			err = e
			util.LogErrorf("internal.Service.infer, parse float64 failed, error: %v", e)
			return
		}
		// 保留1
		prediction = int(p) + 1
	}
	return
}

func (s *Service) GetName() string {
	return s.serviceName
}
