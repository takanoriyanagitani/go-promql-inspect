package qlinspect

import (
	"context"
	"io"
	"strings"

	p2 "github.com/prometheus/prometheus/promql/parser"
)

type Inspector func(p2.Node, []p2.Node) error

func (i Inspector) AsFunc() func(p2.Node, []p2.Node) error { return i }

func (i Inspector) Inspect(node p2.Node) { p2.Inspect(node, i.AsFunc()) }

type ParseOpts p2.Options

var ParseOptsDefault ParseOpts

//nolint:ireturn
func (o ParseOpts) ToParser() p2.Parser { return p2.NewParser(p2.Options(o)) }

type RawSource func(context.Context) (string, error)

type ExprParser func(string) (p2.Expr, error)

type ExprSource func(context.Context) (p2.Expr, error)

func (s RawSource) ToExprSource(conv ExprParser) ExprSource {
	return func(ctx context.Context) (p2.Expr, error) {
		raw, err := s(ctx)
		if nil != err {
			return nil, err
		}

		return conv(raw)
	}
}

type NodeSource func(context.Context) (p2.Node, error)

func (s NodeSource) ToInspector(ctx context.Context, i Inspector) error {
	node, err := s(ctx)
	if nil != err {
		return err
	}

	i.Inspect(node)

	return nil
}

func (e ExprSource) ToNodeSource() NodeSource {
	return func(ctx context.Context) (p2.Node, error) { return e(ctx) }
}

type Parser struct{ p2.Parser }

func (p Parser) AsExprParser() ExprParser { return p.Parser.ParseExpr }

type NodeInspector func(p2.Node) error

func (n NodeInspector) ToInspector() Inspector {
	return func(node p2.Node, _ []p2.Node) error { return n(node) }
}

type IoReader struct{ io.Reader }

func (i IoReader) ToRawSource(limit int64) RawSource {
	return func(_ context.Context) (string, error) {
		var buf strings.Builder
		lmtd := &io.LimitedReader{R: i.Reader, N: limit}
		_, err := io.Copy(&buf, lmtd)
		return buf.String(), err
	}
}
