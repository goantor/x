package x

var (
	_isRelease bool   = false
	_release   string = "release"
	_env       string = "debug"
)

func SetEnv(release string, current func() string) {
	_env = current()
	_release = release
	_isRelease = _env == release
}

func IsRelease() bool {
	return _isRelease
}

func Env() string {
	return _env
}
