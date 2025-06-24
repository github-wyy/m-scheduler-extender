package main

import (
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"net/http"
	"sync"

	"k8s.io/klog/v2"
)

const (
	filterPrefix    = "/scheduler/filter"
	scorePrefix     = "/scheduler/score"
	healthCheckPath = "/healthz"
)

func main() {
	// handler
	http.HandleFunc(healthCheckPath, loggingMiddleware(healthCheckHandler))
	http.HandleFunc(filterPrefix, loggingMiddleware(filterHandler))
	http.HandleFunc(scorePrefix, loggingMiddleware(scoreHandler))

	// web server
	var wg sync.WaitGroup
	wg.Add(2)
	// 启动HTTP服务(8010端口)
	go func() {
		defer wg.Done()
		klog.Info("Starting scheduler extender on :8010")
		if err := http.ListenAndServe(":8010", nil); err != nil {
			klog.Fatalf("Failed to start server: %v", err)
		}
	}()
	// 启动HTTPS服务(8443端口)
	go func() {
		defer wg.Done()
		klog.Info("Starting HTTPS server on :8443")

		// ssl
		certFile := "./ssl/cert.pem"
		keyFile := "./ssl/key.pem"

		tlsConfig := &tls.Config{
			// ClientAuth 表示程序作为服务端时要不要验证客户端证书，
			// InsecureSkipVerify 表示程序作为客户端时要不要跳过服务端证书的验证。
			// ClientAuth: tls.RequireAndVerifyClientCert,
			InsecureSkipVerify: true,
			ClientCAs:          loadKubernetesCA(),
		}
		server := &http.Server{
			Addr:      ":8443",
			TLSConfig: tlsConfig,
		}
		server.Handler = http.DefaultServeMux
		if err := server.ListenAndServeTLS(certFile, keyFile); err != nil {
			klog.Errorf("HTTPS server error: %v", err)
		}
	}()

	wg.Wait()
	klog.Info("All servers stopped")
}

// 加载 Kubernetes CA 证书
func loadKubernetesCA() *x509.CertPool {
	caPath := "/var/run/secrets/kubernetes.io/serviceaccount/ca.crt"
	caCert, err := ioutil.ReadFile(caPath)
	if err != nil {
		klog.Fatalf("Failed to read Kubernetes CA certificate: %v", err)
	}

	caCertPool := x509.NewCertPool()
	if !caCertPool.AppendCertsFromPEM(caCert) {
		klog.Fatalf("Failed to parse Kubernetes CA certificate")
	}

	return caCertPool
}

var loggingMiddleware = func(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// 打印请求基本信息
		klog.Infof("收到请求: %s %s", r.Method, r.URL.Path)
		klog.Infof("远程地址: %s", r.RemoteAddr)

		// 专注于证书和token相关的内容
		klog.Info("证书和token相关信息:")

		// 检查Authorization头信息
		if authHeader := r.Header.Get("Authorization"); authHeader != "" {
			klog.Infof("  Authorization头: %s", authHeader)

			// 提取Bearer Token信息
			if len(authHeader) > 7 && authHeader[:7] == "Bearer " {
				token := authHeader[7:]
				klog.Infof("  Bearer Token: %s... (长度: %d)", token[:min(10, len(token))], len(token))
			}
		} else {
			klog.Info("  Authorization头: 未提供")
		}

		// 检查客户端证书相关的头信息
		certHeaders := []string{"X-Client-Cert", "X-Forwarded-Client-Cert", "Client-Cert"}
		for _, header := range certHeaders {
			if value := r.Header.Get(header); value != "" {
				klog.Infof("  %s: %s... (长度: %d)", header, value[:min(20, len(value))], len(value))
			} else {
				klog.Infof("  %s: 未提供", header)
			}
		}

		// 检查TLS连接状态和客户端证书
		if r.TLS != nil {
			klog.Info("  TLS连接状态:")

			// 检查客户端证书
			if len(r.TLS.PeerCertificates) > 0 {
				klog.Info("  客户端证书信息:")
				for i, cert := range r.TLS.PeerCertificates {
					klog.Infof("    证书 #%d:", i+1)
					klog.Infof("      主题: %s", cert.Subject.CommonName)
					klog.Infof("      颁发者: %s", cert.Issuer.CommonName)
					klog.Infof("      有效期: %s - %s", cert.NotBefore.Format("2006-01-02"), cert.NotAfter.Format("2006-01-02"))
					klog.Infof("      DNS名称: %v", cert.DNSNames)
				}
			} else {
				klog.Info("  客户端证书: 未提供")
			}
		} else {
			klog.Info("  TLS连接状态: 非TLS连接")
		}

		// 调用原始处理器
		next(w, r)
	}
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
