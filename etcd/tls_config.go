/**
 * @Author:  jager
 * @Email:   lhj168os@gmail.com
 * @File:    tls_config
 * @Date:    2021/6/11 2:25 下午
 * @package: etcd
 * @Version: v1.0.0
 *
 * @Description:
 *
 */

package etcd

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/jageros/hawos/log"
	"io/ioutil"
	"net"
	"net/http"
	"time"
)

var etcdCertPath = "etcdclient/kube-etcd.pem"
var etcdKeyPath = "etcdclient/kube-etcd-key.pem"
var etcdCAPath = "etcdclient/kube-ca.pem"

var _tlsConfig *tls.Config

func TLSConfig() (*tls.Config, error) {
	if _tlsConfig != nil {
		return _tlsConfig, nil
	}

	// load cert
	cert, err := tls.LoadX509KeyPair(etcdCertPath, etcdKeyPath)
	if err != nil {
		return nil, err
	}

	// load root ca
	caData, err := ioutil.ReadFile(etcdCAPath)
	if err != nil {
		return nil, err
	}
	pool := x509.NewCertPool()
	pool.AppendCertsFromPEM(caData)

	_tlsConfig = &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      pool,
	}
	return _tlsConfig, nil
}

func TlsTransport() (tr *http.Transport) {
	tlsConfig, err := TLSConfig()
	if err != nil {
		log.Debugf("Create tls error, error=\"%s\"", err.Error())
		return
	}
	tr = &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		TLSHandshakeTimeout: 10 * time.Second,
		TLSClientConfig:     tlsConfig,
	}
	return
}
