// Package grpc provides a gRPC options
package grpc

import (
	"crypto/tls"
	"net"
	"time"

	"github.com/micro/go-micro/client"
	"golang.org/x/net/context"
	"google.golang.org/grpc"
)

type codecsKey struct{}
type tlsAuth struct{}
type dialerKey struct{}

// gRPC Codec to be used to encode/decode requests for a given content type
func Codec(contentType string, c grpc.Codec) client.Option {
	return func(o *client.Options) {
		codecs := make(map[string]grpc.Codec)
		if o.Context == nil {
			o.Context = context.Background()
		}
		if v := o.Context.Value(codecsKey{}); v != nil {
			codecs = v.(map[string]grpc.Codec)
		}
		codecs[contentType] = c
		o.Context = context.WithValue(o.Context, codecsKey{}, codecs)
	}
}

// AuthTLS should be used to setup a secure authentication using TLS
func AuthTLS(t *tls.Config) client.Option {
	return func(o *client.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, tlsAuth{}, t)
	}
}

func Dialer(d func(addr string, timeout time.Duration) (net.Conn, error)) client.Option {
	return func(o *client.Options) {
		if o.Context == nil {
			o.Context = context.Background()
		}
		o.Context = context.WithValue(o.Context, dialerKey{}, d)
	}
}

func ProxyDialer(proxyAddr string) client.Option {
	return Dialer(func(addr string, timeout time.Duration) (net.Conn, error) {
		return net.DialTimeout("tcp", proxyAddr, timeout)
	})
}
