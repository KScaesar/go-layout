//go:build !intg

package utility

// UpDocker 由於 TestMain 在 unit test 也會呼叫 UpDocker
// 為了避免啟動多餘的容器, 設計一個無作用的函數
func UpDocker(removeData bool, services []DockerService) (DownDocker func() error) {
	return func() error { return nil }
}
