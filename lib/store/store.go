package store

type Store interface {
  AddDevice(name string) error
  ListDevices() ([]string, error)
  DeviceExists(name string) bool
  RemoveDevice(name string) error
}
