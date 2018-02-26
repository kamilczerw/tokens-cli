package store

import (
  "path/filepath"
  "os"
  "errors"
  "io/ioutil"
)

func NewFileStore() (*FileStore, error) {
  fileStore := new(FileStore)
  err := fileStore.init()

  if err != nil {
    return nil, err
  }

  return fileStore, nil
}

func GetPath(path string, defaultPath string) string {
  if path != "" && filepath.IsAbs(string(path)) {
    return path
  }

  return defaultPath
}

func (store *FileStore) init() error {
  home := os.Getenv("HOME")
  defaultDataPath := filepath.Join(home, ".tokens")

  store.DataPath = GetPath(os.Getenv("TOKENS_DATA_PATH"), defaultDataPath)

  return os.MkdirAll(store.DataPath, 0700)
}

type FileStore struct {
  DataPath string
}

func (store *FileStore) AddDevice(name string) error {
  if store.DeviceExists(name) {
    return errors.New("device already exists")
  }

  _, err := os.Create(filepath.Join(store.DataPath, name))
  if err != nil {
    return err
  }

  return nil
}

func (store *FileStore) ListDevices() ([]string, error) {
  files, err := ioutil.ReadDir(store.DataPath)
  if err != nil {
    return nil, err
  }

  var devices []string = nil
  for _, file := range files {
    // TODO: Verify if the device has secret stored int keyring
    devices = append(devices, file.Name())
  }

  return devices, nil
}

func (store *FileStore) DeviceExists(name string) bool {
  _, err := os.Stat(filepath.Join(store.DataPath, name))

  return err == nil
}

func (store *FileStore) RemoveDevice(name string) error {
  if !store.DeviceExists(name) {
    return errors.New("device does not exist")
  }

  err := os.Remove(filepath.Join(store.DataPath, name))
  if err != nil {
    return err
  }

  return nil
}
