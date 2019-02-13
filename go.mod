module chuck

require (
	chuck/cmds v0.0.0
	chuck/handlers v0.0.0
	chuck/storage v0.0.0
	chuck/utils v0.0.0
	github.com/gavrilaf/chuck v0.0.0-20190212120714-73f3027f5761 // indirect

	github.com/mitchellh/cli v1.0.0
	github.com/spf13/afero v1.2.0
)

replace chuck/cmds => ./cmds

replace chuck/storage => ./storage

replace chuck/handlers => ./handlers

replace chuck/utils => ./utils
