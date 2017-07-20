package main

import (
	"crypto/tls"
	"crypto/x509"
	"flag"
	"fmt"
	"io/ioutil"
	"log"

	redis "github.com/go-redis/redis"
)

func main() {
	address := flag.String("address", "", "The address:port of the redis server")
	certFile := flag.String("cert", "", "")
	keyFile := flag.String("key", "", "")
	caCert := flag.String("cacert", "", "")
	flag.Parse()

	opts, err := redis.ParseURL(fmt.Sprintf("rediss://%s/0", *address))
	if err != nil {
		log.Fatal(err)
	}

	cert, err := buildClientCert(*certFile, *keyFile)
	if err != nil {
		log.Fatal(err)
	}
	opts.TLSConfig.Certificates = []tls.Certificate{cert}

	if *caCert != "" {
		caPool, err := buildCAPool(*caCert)
		if err != nil {
			log.Fatal(err)
		}
		opts.TLSConfig.RootCAs = caPool
	}

	client := redis.NewClient(opts)

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)
}

func buildClientCert(certFile, keyFile string) (tls.Certificate, error) {
	return tls.LoadX509KeyPair(certFile, keyFile)
}

func buildCAPool(caCertFile string) (*x509.CertPool, error) {
	cacert, err := ioutil.ReadFile(caCertFile)
	if err != nil {
		return nil, err
	}
	pool := x509.NewCertPool()
	ok := pool.AppendCertsFromPEM(cacert)
	if !ok {
		return nil, fmt.Errorf("Failed to parse CA cert")
	}
	return pool, nil
}
