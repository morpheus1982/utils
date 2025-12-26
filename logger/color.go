package logger

import (
	"go.uber.org/zap/buffer"
	"go.uber.org/zap/zapcore"
)

// ANSI 颜色代码
const (
	colorReset  = "\033[0m"
	colorRed    = "\033[31m"
	colorYellow = "\033[33m"
	colorBlue   = "\033[34m"
	colorGray   = "\033[90m"
	colorCyan   = "\033[96m"
	colorGreen  = "\033[32m"
	colorPurple = "\033[35m"
)

// colorLevelEncoder 为不同日志等级添加颜色（不重置颜色，让颜色延续到整行）
func colorLevelEncoder(l zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
	// 直接输出等级文本，不添加颜色，颜色由 EncodeEntry 统一处理
	enc.AppendString(l.CapitalString())
}

// coloredConsoleEncoder 包装 zapcore.Encoder，为整行添加颜色
type coloredConsoleEncoder struct {
	zapcore.Encoder
	level zapcore.Level
}

// Clone 实现 Encoder 接口
func (e *coloredConsoleEncoder) Clone() zapcore.Encoder {
	return &coloredConsoleEncoder{
		Encoder: e.Encoder.Clone(),
		level:   e.level,
	}
}

// EncodeEntry 实现 Encoder 接口，为整行日志添加颜色
func (e *coloredConsoleEncoder) EncodeEntry(entry zapcore.Entry, fields []zapcore.Field) (*buffer.Buffer, error) {
	e.level = entry.Level
	buf, err := e.Encoder.EncodeEntry(entry, fields)
	if err != nil {
		return nil, err
	}

	// 根据日志等级为整行添加颜色
	coloredBuf := buffer.NewPool().Get()
	var color string
	switch entry.Level {
	case zapcore.DebugLevel:
		color = colorGray
	case zapcore.InfoLevel:
		color = colorCyan
	case zapcore.WarnLevel:
		color = colorYellow
	case zapcore.ErrorLevel:
		color = colorRed
	case zapcore.DPanicLevel, zapcore.PanicLevel:
		color = colorPurple
	case zapcore.FatalLevel:
		color = colorRed
	default:
		color = colorReset
	}

	coloredBuf.AppendString(color)
	coloredBuf.Write(buf.Bytes())
	// 移除末尾的换行符（如果有）
	if coloredBuf.Len() > 0 && coloredBuf.Bytes()[coloredBuf.Len()-1] == '\n' {
		coloredBuf.TrimNewline()
		coloredBuf.AppendString(colorReset)
		coloredBuf.AppendByte('\n')
	} else {
		coloredBuf.AppendString(colorReset)
	}

	buf.Free()
	return coloredBuf, nil
}

// newColoredConsoleEncoder 创建一个带颜色的控制台编码器
func newColoredConsoleEncoder(cfg zapcore.EncoderConfig) zapcore.Encoder {
	// 为 level 设置颜色编码器
	cfg.EncodeLevel = colorLevelEncoder
	return &coloredConsoleEncoder{
		Encoder: zapcore.NewConsoleEncoder(cfg),
	}
}