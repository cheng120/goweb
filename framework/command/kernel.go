package command

import "goweb/framework/cobra"

// AddKernelCommands will add all command/* to root command
func AddKernelCommands(root *cobra.Command) {
	// root.AddCommand(DemoCommand)
	root.AddCommand(initEnvCommand())
	root.AddCommand(initEnvCommand())
	//root.AddCommand(deployCommand)
	// config 命令
	root.AddCommand(initConfigCommand())
	//// cron
	root.AddCommand(initCronCommand())
	//// cmd
	//cmdCommand.AddCommand(cmdListCommand)
	//cmdCommand.AddCommand(cmdCreateCommand)
	root.AddCommand(initCmdCommand())
	//
	//// build
	// build 命令
	root.AddCommand(initBuildCommand())
	// go build
	root.AddCommand(goCommand)
	// npm build
	root.AddCommand(npmCommand)
	//
	//// app
	root.AddCommand(initAppCommand())
	//
	// dev
	root.AddCommand(initDevCommand())
	//
	//// middleware
	//middlewareCommand.AddCommand(middlewareAllCommand)
	//middlewareCommand.AddCommand(middlewareAddCommand)
	//middlewareCommand.AddCommand(middlewareRemoveCommand)
	root.AddCommand(initMiddlewareCommand())
	//
	//// swagger
	//swagger.IndexCommand.AddCommand(swagger.InitServeCommand())
	//swagger.IndexCommand.AddCommand(swagger.GenCommand)
	//root.AddCommand(swagger.IndexCommand)
	//
	//// provider
	root.AddCommand(initProviderCommand())
	// providerCommand.AddCommand(providerListCommand)
	// providerCommand.AddCommand(providerCreateCommand)
	// root.AddCommand(providerCommand)
	//
	// new
	root.AddCommand(initNewCommand())
}