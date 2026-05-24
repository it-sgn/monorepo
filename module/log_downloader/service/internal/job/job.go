package job111

import "github.com/google/wire"

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewExampleJob, NewDownloadJob)

type JobFunc func()

var DefaultJobs map[string]JobFunc

// type DownloadFunc func()

// var DefaultDownLoads map[string]DownloadFunc
