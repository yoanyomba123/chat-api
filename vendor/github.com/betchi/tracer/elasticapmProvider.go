package tracer

import (
	"context"
	"net/http"

	elasticapm "github.com/elastic/apm-agent-go"
)

var elasticapmTracer *elasticapm.Tracer

type elasticapmProvider struct {
	logger elasticapm.Logger
}

func (ep *elasticapmProvider) NewTracer(config *Config) error {
	if ep.logger != nil {
		elasticapmTracer.SetLogger(ep.logger)
	}

	elasticapmTracer, err := elasticapm.NewTracer(config.ServiceName, config.ServiceVersion)
	if err != nil {
		if ep.logger != nil {
			ep.logger.Errorf(err.Error())
		}
		return err
	}

	elasticapmTracer.SetCaptureBody(elasticapm.CaptureBodyAll)
	return nil
}

func (ep *elasticapmProvider) StartTransaction(ctx context.Context, name, transactionType string, opts ...StartTransactionOption) (context.Context, interface{}) {
	if elasticapmTracer == nil {
		return ctx, nil
	}

	transaction := elasticapmTracer.StartTransaction(name, transactionType)
	ctxWithTx := elasticapm.ContextWithTransaction(ctx, transaction)
	ctx = context.WithValue(ctxWithTx, CtxTracerSpan, transaction)
	return ctx, transaction
}

func (ep *elasticapmProvider) StartSpan(ctx context.Context, name, spanType string) interface{} {
	if elasticapmTracer == nil {
		return nil
	}

	span, _ := elasticapm.StartSpan(ctx, name, spanType)
	return span
}

func (ep *elasticapmProvider) InjectHTTPRequest(span interface{}, req *http.Request) {
	if span == nil {
		return
	}
}

func (ep *elasticapmProvider) SetTag(span interface{}, key string, value interface{}) {
	// transaction := ep.ctx.Value(CtxTracerSpan)
	// if transaction != nil {
	// 	txCtx := transaction.(*elasticapm.Transaction).Context
	// 	// txCtx.SetTag(key, fmt.Sprintf("%v", value))
	// 	txCtx.SetCustom(key, value)
	// }
}

func (ep *elasticapmProvider) SetHTTPStatusCode(span interface{}, statusCode int) {
	// transaction := ep.ctx.Value(CtxTracerSpan)
	// if transaction != nil {
	// 	txCtx := transaction.(*elasticapm.Transaction).Context
	// 	txCtx.SetHTTPStatusCode(statusCode)
	// }
}

func (ep *elasticapmProvider) SetError(span interface{}, err error) {
	// TODO
}

func (ep *elasticapmProvider) Finish(span interface{}) {
	if span != nil {
		span.(*elasticapm.Span).End()
	}
}

func (ep *elasticapmProvider) CloseTransaction(ctx context.Context) {
	transaction := ctx.Value(CtxTracerSpan)
	if transaction != nil {
		transaction.(*elasticapm.Transaction).End()
	}

	elasticapmTracer.Flush(nil)
}

func (ep *elasticapmProvider) Close() {
}
