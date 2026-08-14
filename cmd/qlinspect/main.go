package main

import (
	"context"
	"io"
	"log/slog"
	"os"
	"reflect"
	"strconv"

	p2 "github.com/prometheus/prometheus/promql/parser"
	pi "github.com/takanoriyanagitani/go-promql-inspect"
)

var smaxQl string = os.Getenv("ENV_PROMQL_MAX_SIZE")

var irdr io.Reader = os.Stdin

var popts pi.ParseOpts = pi.ParseOptsDefault
var parser pi.Parser = pi.Parser{Parser: popts.ToParser()}
var eparser pi.ExprParser = parser.AsExprParser()

var node2log pi.NodeInspector = func(node p2.Node) error {
	if nil == node {
		return nil
	}

	switch t := node.(type) {
	case *p2.Call:
		slog.Info("func got", "name", t.Func.Name)
		return nil
	case *p2.MatrixSelector:
		slog.Info("matrix got", "duration", t.Range)
		return nil
	case *p2.VectorSelector:
		slog.Info("vector got", "name", t.Name)
		for _, lmat := range t.LabelMatchers {
			slog.Info(
				"  matcher got",
				"name", lmat.Name,
				"value", lmat.Value,
				"type", lmat.Type,
			)
		}
		return nil
	case *p2.AggregateExpr:
		slog.Info("aggregate got", "short", t.ShortString())
		for _, grp := range t.Grouping {
			slog.Info(
				"  group got",
				"group", grp,
			)
		}
		return nil
	case *p2.BinaryExpr:
		slog.Info("binary expr got", "short", t.ShortString())
		slog.Info("binary expr", "lhs", t.LHS.String())
		slog.Info("binary expr", "rhs", t.RHS.String())
		return nil
	case *p2.ParenExpr:
		slog.Info("paren expr got", "long", t.String())
		return nil
	case *p2.NumberLiteral:
		slog.Info("number got", "val", t.Val, "isDuration", t.Duration)
		return nil
	default:
	}

	var rtyp reflect.Type = reflect.TypeOf(node)

	slog.Info(
		"node got", "node", rtyp.String(),
	)
	return nil
}

var insp pi.Inspector = node2log.ToInspector()

var sub func(context.Context) error = func(ctx context.Context) error {
	slog.Info("setting up", "ENV_PROMQL_MAX_SIZE", smaxQl)

	imaxQl, err := strconv.Atoi(smaxQl)
	if nil != err {
		return err
	}

	var rsrc pi.RawSource = pi.IoReader{Reader: irdr}.ToRawSource(int64(imaxQl))
	var esrc pi.ExprSource = rsrc.ToExprSource(eparser)
	var nsrc pi.NodeSource = esrc.ToNodeSource()

	slog.Info("configured", "promql max size", imaxQl)

	return nsrc.ToInspector(ctx, insp)
}

func main() {
	err := sub(context.Background())
	if nil != err {
		slog.Error("error got", "error", err)
	}
}
