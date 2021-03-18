// Copyright 2016-2021 The Libsacloud Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package otel

import (
	"github.com/sacloud/libsacloud/v2"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

type config struct {
	Tracer           trace.Tracer
	TracerProvider   trace.TracerProvider
	SpanStartOptions []trace.SpanOption
}

type Option interface {
	Apply(*config)
}

type OptionFunc func(*config)

func (o OptionFunc) Apply(c *config) {
	o(c)
}

func newConfig(opts ...Option) *config {
	c := &config{
		TracerProvider: otel.GetTracerProvider(),
	}
	for _, opt := range opts {
		opt.Apply(c)
	}

	c.Tracer = c.TracerProvider.Tracer(
		"github.com/sacloud/libsacloud",
		trace.WithInstrumentationVersion(libsacloud.Version),
	)
	return c
}

func WithTracerProvider(provider trace.TracerProvider) Option {
	return OptionFunc(func(cfg *config) {
		cfg.TracerProvider = provider
	})
}

func WithSpanOptions(opts ...trace.SpanOption) Option {
	return OptionFunc(func(c *config) {
		c.SpanStartOptions = append(c.SpanStartOptions, opts...)
	})
}
