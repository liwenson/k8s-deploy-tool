package config

var (
	Keepdirectory = true
)

// 参数判断
func IsArgeNull(args Args) {

	if args.Pid == "" {
		panic("project ID cannot be empty")
	}

	if args.Token == "" {
		panic("Encryption token cannot be empty")
	}

	if args.Server == "" {
		panic("server cannot be empty")
	}

	if args.Out != "" {
		// 指定了输出路径就不保留一级目录
		abc := args.Out[len(args.Out)-1:]
		if abc == "/" {
			Keepdirectory = false
		}

	}

}
