package config

type Config struct {
	Port         int    `yaml:"port"`
	UploadDir    string `yaml:"upload_dir"`
	DownloadDir  string `yaml:"download_dir"`
	LogDir       string `yaml:"log_dir"`
	DatabasePath string `yaml:"database_path"`
	MaxFileSize  int64  `yaml:"max_file_size"`
}

func Load() *Config {
	return &Config{
		Port:         8080,
		UploadDir:    "./uploads",   // 上传目录配置
		DownloadDir:  "./downloads", // 下载目录配置
		LogDir:       "./logs",      // 日志目录配置
		DatabasePath: "./file-service.db",
		MaxFileSize:  5 << 20, // 5MB
	}
}
