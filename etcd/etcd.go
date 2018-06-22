package etcd

import (
  "fmt"
  "io/ioutil"
  "net"
  "net/http"
  "net/url"
  "os"

  "github.com/coreos/etcd/etcdserver"
  "github.com/coreos/etcd/etcdserver/api/v2http"
  "github.com/coreos/etcd/pkg/types"
)

const (
  memberName  = "simple"
  clusterName = "simple-cluster"
  tempPrefix  = "simple-etcd-"

  // No peer URL exists but etcd doesn't allow the value to be empty.
  peerURL    = "http://localhost:0"
  clusterCfg = memberName + "=" + peerURL
)

// SimpleEtcd provides a single node etcd server.
type SimpleEtcd struct {
  Port     int
  listener net.Listener
  server   *etcdserver.EtcdServer
  dataDir  string
}

func NewSimpleEtcd() (*SimpleEtcd, error) {
  var err error
  se := &SimpleEtcd{}
  se.listener, err = net.Listen("tcp", ":0")
  if err != nil {
    return nil, err
  }

  se.Port = se.listener.Addr().(*net.TCPAddr).Port
  clientURLs, err := interfaceURLs(se.Port)
  if err != nil {
    se.Destroy()
    return nil, err
  }

  se.dataDir, err = ioutil.TempDir("", tempPrefix)
  if err != nil {
    se.Destroy()
    return nil, err
  }

  peerURLs, err := types.NewURLs([]string{peerURL})
  if err != nil {
    se.Destroy()
    return nil, err
  }

  cfg := &etcdserver.ServerConfig{
    Name:       memberName,
    ClientURLs: clientURLs,
    PeerURLs:   peerURLs,
    DataDir:    se.dataDir,
    InitialPeerURLsMap: types.URLsMap{
      memberName: peerURLs,
    },
    NewCluster:    true,
    TickMs:        100,
    ElectionTicks: 10,
  }

  se.server, err = etcdserver.NewServer(cfg)
  if err != nil {
    return nil, err
  }

  se.server.Start()
  go http.Serve(se.listener,
    v2http.NewClientHandler(se.server, cfg.ReqTimeout()))

  return se, nil
}

func (se *SimpleEtcd) Destroy() {
  if se.listener != nil {
    if err := se.listener.Close(); err != nil {
      plog.Errorf("Error closing etcd listener: %v", err)
    }
  }

  if se.server != nil {
    se.server.Stop()
  }

  if se.dataDir != "" {
    if err := os.RemoveAll(se.dataDir); err != nil {
      plog.Errorf("Error removing etcd data dir: %v", err)
    }
  }
}

// Generate all publishable URLs for a given HTTP port.
func interfaceURLs(port int) (types.URLs, error) {
  allAddrs, err := net.InterfaceAddrs()
  if err != nil {
    return []url.URL{}, err
  }

  var allURLs types.URLs
  for _, a := range allAddrs {
    ip, ok := a.(*net.IPNet)
    if !ok || !ip.IP.IsGlobalUnicast() {
      continue
    }

    tcp := net.TCPAddr{
      IP:   ip.IP,
      Port: port,
    }

    u := url.URL{
      Scheme: "http",
      Host:   tcp.String(),
    }
    allURLs = append(allURLs, u)
  }

  if len(allAddrs) == 0 {
    return []url.URL{}, fmt.Errorf("no publishable addresses")
  }

  return allURLs, nil
}
