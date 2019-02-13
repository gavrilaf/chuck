module chuck/cmds

require (
	chuck/handlers v0.0.0
	chuck/utils v0.0.0
	github.com/spf13/afero v1.2.0
	gopkg.in/elazarl/goproxy.v1 v1.0.0-20180725130230-947c36da3153
)

replace chuck/utils => ../utils

replace chuck/storage => ../storage

replace chuck/handlers => ../handlers
