module chuck/handlers

require (
	chuck/storage v0.0.0
	chuck/utils v0.0.0
	github.com/elazarl/goproxy v0.0.0-20181111060418-2ce16c963a8a // indirect
	github.com/mitchellh/cli v1.0.0
	github.com/onsi/ginkgo v1.7.0
	github.com/onsi/gomega v1.5.0
	github.com/spf13/afero v1.2.0
	gopkg.in/elazarl/goproxy.v1 v1.0.0-20180725130230-947c36da3153

)

replace chuck/utils => ../utils

replace chuck/storage => ../storage
