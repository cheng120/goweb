package console

import (
	"goweb/app/console/command/demo"
	"goweb/framework"
	"goweb/framework/cobra"
	"goweb/framework/command"
)

// RunCommand is command
func RunCommand(container framework.Container) error {
	var rootCmd = &cobra.Command{
		Use:   "main",
		Short: "main 命令",
		Long:  "main 框架提供的命令行工具，使用这个命令行工具能很方便执行框架自带命令，也能很方便编写业务命令",
		RunE: func (cmd *cobra.Command,args []string) error  {
			cmd.InitDefaultHelpFlag()
			return cmd.Help()
		},
		// 不需要出现cobra默认的completion子命令
		CompletionOptions: cobra.CompletionOptions{DisableDefaultCmd: true},
	}
	// 为根Command设置服务容器
	rootCmd.SetContainer(container)
	// 绑定框架的命令
	command.AddKernelCommands(rootCmd)
	// 绑定业务的命令
	AddAppCommand(rootCmd)
	// 执行RootCommand
	// rootCmd.AddCronCommand("* * * * *", command.DemoCommand)

	return rootCmd.Execute()
}


func AddAppCommand(cmd *cobra.Command) {
	cmd.AddCommand(demo.FooCommand)
	// 每秒调用一次Foo命令
	//rootCmd.AddCronCommand("* * * * * *", demo.FooCommand)

	// 启动一个分布式任务调度，调度的服务名称为init_func_for_test，每个节点每5s调用一次Foo命令，抢占到了调度任务的节点将抢占锁持续挂载2s才释放
	//rootCmd.AddDistributedCronCommand("foo_func_for_test", "*/5 * * * * *", demo.FooCommand, 2*time.Second)
}