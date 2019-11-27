package config

import (
	"crypto/tls"
	"fmt"
	"strconv"

	"github.com/go-redis/redis"
)

// GetRedis function
func GetRedis(redisHost, redisTLS, redisPassword, redisPort string) (*redis.Client, error) {
	//Transport Layer Security config,
	// If InsecureSkipVerify is true, TLS accepts any certificate
	// presented by the server and any host name in that certificate.
	// https://godoc.org/crypto/tls#Config

	tlsSecured, err := strconv.ParseBool(redisTLS)
	if err != nil {
		return nil, err
	}

	var conf *tls.Config

	// force checking for unsecured aws redis
	if tlsSecured {
		conf = &tls.Config{
			InsecureSkipVerify: tlsSecured,
		}
	} else {
		conf = nil
	}

	cl := redis.NewClient(&redis.Options{
		Addr:      fmt.Sprintf("%v:%v", redisHost, redisPort),
		Password:  redisPassword,
		DB:        0, // use default DB
		TLSConfig: conf,
	})

	return cl, nil
}
