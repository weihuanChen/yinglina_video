package configs

import (
	"flag"
	"fmt"
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
	"os"
	"path/filepath"
	"reflect"
)

type app struct {
	Version      string `mapstructure:"version"`
	Name         string `mapstructure:"name"`
	IdleTimeOut  int    `mapstructure:"idle_time_out"`
	ReadTimeOut  int    `mapstructure:"read_time_out"`
	WriteTimeOut int    `mapstructure:"write_time_out"`
	Port         string `mapstructure:"port"`
	Env          string `mapstructure:"env"`
	Desc         string `mapstructure:"desc"`
	CfgFile      string `mapstructure:"-" ignore:"true"` // 记录实际使用了哪个配置文件, 这是运行时配置, 不从 YAML 加载; 且忽略变更时类型检查
}

type logger struct {
	LogFilePath     string `mapstructure:"logFilePath"`
	LogFileName     string `mapstructure:"logFileName"`
	LogTimestampFmt string `mapstructure:"logTimestampFmt"`
	LogMaxAge       int    `mapstructure:"logMaxAge"`
	LogRotationTime int    `mapstructure:"logRotationTime"`
	LogLevel        string `mapstructure:"logLevel"`
}

type db struct {
	Name string `mapstructure:"name"`
	Host string `mapstructure:"host"`
	Port string `mapstructure:"port"`
	User string `mapstructure:"user"`
	Psw  string `mapstructure:"psw"`
}

type redis struct {
	Host string `yaml:"host"`
	Port string `yaml:"port"`
	Db   string `yaml:"db"`
	Psw  string `yaml:"psw"`
}

type AppConfig struct {
	App    app    `mapstructure:"app"`
	Logger logger `mapstructure:"logger"`
	Db     db     `mapstructure:"db"`
	Redis  redis  `mapstructure:"redis"`
}

// CfgPath 配置文件路径
const (
	CfgPath = "configs/app.yaml"
)

// Cfg 配置信息
var (
	Cfg    AppConfig
	oldCfg AppConfig // 用于存储配置变更前的副本信息
)

// deepCopy 使用反射实现深拷贝
func deepCopy(src interface{}) interface{} {
	srcValue := reflect.ValueOf(src)
	if srcValue.Kind() == reflect.Ptr {
		srcValue = srcValue.Elem()
	}

	// 创建新的目标值
	dstValue := reflect.New(srcValue.Type()).Elem()
	deepCopyValue(srcValue, dstValue)

	return dstValue.Interface()
}

// deepCopyValue 递归复制值
func deepCopyValue(src, dst reflect.Value) {
	switch src.Kind() {
	case reflect.Struct:
		for i := 0; i < src.NumField(); i++ {
			deepCopyValue(src.Field(i), dst.Field(i))
		}
	case reflect.Slice, reflect.Array:
		if src.IsNil() {
			return
		}
		dst.Set(reflect.MakeSlice(src.Type(), src.Len(), src.Cap()))
		for i := 0; i < src.Len(); i++ {
			deepCopyValue(src.Index(i), dst.Index(i))
		}
	case reflect.Map:
		if src.IsNil() {
			return
		}
		dst.Set(reflect.MakeMap(src.Type()))
		for _, key := range src.MapKeys() {
			value := src.MapIndex(key)
			dst.SetMapIndex(key, value)
		}
	default:
		if dst.CanSet() {
			dst.Set(src)
		}
	}
}

// compareAndPrintChanges 通用的配置比较函数
func compareAndPrintChanges(old, new interface{}, path string) {
	oldVal := reflect.ValueOf(old)
	newVal := reflect.ValueOf(new)

	// 如果是指针，获取其指向的值
	if oldVal.Kind() == reflect.Ptr {
		oldVal = oldVal.Elem()
	}
	if newVal.Kind() == reflect.Ptr {
		newVal = newVal.Elem()
	}

	oldType := oldVal.Type()

	// 遍历所有字段
	for i := 0; i < oldVal.NumField(); i++ {
		field := oldType.Field(i)
		oldField := oldVal.Field(i)
		newField := newVal.Field(i)

		// 构建字段路径
		fieldPath := path
		if path != "" {
			fieldPath += "."
		}
		fieldPath += field.Name

		// 根据字段类型进行比较
		switch oldField.Kind() {
		case reflect.Struct:
			// 递归比较结构体
			compareAndPrintChanges(oldField.Interface(), newField.Interface(), fieldPath)

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			oldInt := oldField.Int()
			newInt := newField.Int()
			if oldInt != newInt {
				fmt.Printf("配置项 [%s] 发生变化: %d -> %d\n", fieldPath, oldInt, newInt)
			}

		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			oldUint := oldField.Uint()
			newUint := newField.Uint()
			if oldUint != newUint {
				fmt.Printf("配置项 [%s] 发生变化: %d -> %d\n", fieldPath, oldUint, newUint)
			}

		case reflect.Float32, reflect.Float64:
			oldFloat := oldField.Float()
			newFloat := newField.Float()
			if oldFloat != newFloat {
				fmt.Printf("配置项 [%s] 发生变化: %f -> %f\n", fieldPath, oldFloat, newFloat)
			}

		case reflect.String:
			oldStr := oldField.String()
			newStr := newField.String()
			if oldStr != newStr {
				fmt.Printf("配置项 [%s] 发生变化: %q -> %q\n", fieldPath, oldStr, newStr)
			}

		case reflect.Bool:
			oldBool := oldField.Bool()
			newBool := newField.Bool()
			if oldBool != newBool {
				fmt.Printf("配置项 [%s] 发生变化: %v -> %v\n", fieldPath, oldBool, newBool)
			}

		case reflect.Slice, reflect.Array:
			if !reflect.DeepEqual(oldField.Interface(), newField.Interface()) {
				fmt.Printf("配置项 [%s] 发生变化: %v -> %v\n", fieldPath, oldField.Interface(), newField.Interface())
			}

		case reflect.Map:
			if !reflect.DeepEqual(oldField.Interface(), newField.Interface()) {
				fmt.Printf("配置项 [%s] 发生变化: %v -> %v\n", fieldPath, oldField.Interface(), newField.Interface())
			}

		default:
			if !reflect.DeepEqual(oldField.Interface(), newField.Interface()) {
				fmt.Printf("配置项 [%s] 发生变化: %v -> %v\n", fieldPath, oldField.Interface(), newField.Interface())
			}
		}
	}
}

// validateConfig 校验新配置是否符合 AppConfig 的类型
func validateConfig(newCfg interface{}, referenceCfg interface{}) error {
	newVal := reflect.ValueOf(newCfg)
	refVal := reflect.ValueOf(referenceCfg)

	if newVal.Kind() == reflect.Ptr {
		newVal = newVal.Elem()
	}
	if refVal.Kind() == reflect.Ptr {
		refVal = refVal.Elem()
	}

	for i := 0; i < refVal.NumField(); i++ {
		field := refVal.Type().Field(i)
		newField := newVal.Field(i)
		refField := refVal.Field(i)

		// 跳过带有 ignore 标签的字段
		if field.Tag.Get("ignore") == "true" {
			continue
		}

		// 如果新配置项字段是零值，可能是因为类型不匹配导致未正确解析
		if newField.IsZero() && !refField.IsZero() {
			return fmt.Errorf("配置项 [%s] 类型不匹配或值解析失败", field.Name)
		}

		if newField.Kind() != refField.Kind() {
			return fmt.Errorf("配置项 [%s] 类型不匹配: 期望 %s, 但得到 %s", field.Name, refField.Kind(), newField.Kind())
		}

		// 若字段是结构体，递归校验内部字段类型
		if refField.Kind() == reflect.Struct {
			if err := validateConfig(newField.Interface(), refField.Interface()); err != nil {
				return err
			}
		}
	}
	return nil
}

// reloadCfg 重载配置
func reloadCfg(v *viper.Viper) error {
	return v.Unmarshal(&Cfg)
}

// InitCfg 初始化配置
func InitCfg() {
	var cfgFile string
	flag.StringVar(&cfgFile, "config", "", "配置文件")
	flag.Parse()

	if cfgFile == "" {
		cfgFile, _ = os.Getwd()
		cfgFile = filepath.Join(cfgFile, CfgPath)
	}

	if cfgFile == "" {
		panic("配置文件不存在!")
	}

	v := viper.New()
	v.SetConfigFile(cfgFile)
	v.SetConfigType("yaml")
	if err := v.ReadInConfig(); err != nil {
		panic(fmt.Errorf("配置解析失败: %w", err))
	}

	if err := reloadCfg(v); err != nil {
		panic(fmt.Errorf("配置加载失败: %w", err))
	}

	oldCfg = deepCopy(Cfg).(AppConfig)

	// 观察配置
	v.WatchConfig()
	v.OnConfigChange(func(e fsnotify.Event) {
		fmt.Printf("\n配置文件发生改变: %s\n", e.Name)

		tempCfg := AppConfig{}
		if err := v.Unmarshal(&tempCfg); err != nil {
			fmt.Println("配置重载失败:", err)
			//return
		}

		// 校验新配置类型
		if err := validateConfig(tempCfg, Cfg); err != nil {
			fmt.Println("变更配置类型不匹配，停止变更:", err)
			return
		}

		// 重新加载配置
		if err := reloadCfg(v); err != nil {
			fmt.Println("配置重载失败:", err)
			return
		}

		// 比较并打印变化
		compareAndPrintChanges(oldCfg, Cfg, "App")
		oldCfg = deepCopy(Cfg).(AppConfig)
	})

	Cfg.App.CfgFile = cfgFile
}
