package discovery

import (
	"context"
	"github.com/jageros/hawox/logx"
	"github.com/jageros/hawox/registry"

	"google.golang.org/grpc/resolver"
)

const Name = "discovery"

// Option is builder option.
type Option func(o *builder)

type builder struct {
	ctx        context.Context
	discoverer registry.Discovery
}

// NewBuilder creates a builder which is used to factory registry resolvers.
func NewBuilder(ctx context.Context, d registry.Discovery, opts ...Option) resolver.Builder {
	b := &builder{
		ctx:        ctx,
		discoverer: d,
	}
	for _, o := range opts {
		o(b)
	}
	return b
}

func (d *builder) Build(target resolver.Target, cc resolver.ClientConn, opts resolver.BuildOptions) (resolver.Resolver, error) {
	logx.Info().Str("target", target.Endpoint).Msg("Build Resolver")

	w, err := d.discoverer.Watch(d.ctx, target.Endpoint)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithCancel(d.ctx)
	r := &discoveryResolver{
		w:      w,
		cc:     cc,
		ctx:    ctx,
		cancel: cancel,
	}
	go r.watch()
	return r, nil
}

func (d *builder) Scheme() string {
	return Name
}
