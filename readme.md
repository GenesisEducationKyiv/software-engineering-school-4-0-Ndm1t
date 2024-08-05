# Alerts that should be configured for services
## Basically all the fails should be alerted either those are http request fails or message consume fails
### For http requests
Those are the metrics like request_failed{path="/somepath"} here is the code snippet of collecting such metrics
```
func (s *Server) metricsMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		s := fmt.Sprintf(`requests_total{path=%q}`, ctx.Request.URL.Path)
		metrics.GetOrCreateCounter(s).Inc()
		ctx.Next()
		statusCode := ctx.Writer.Status()
		if statusCode >= 400 {
			s = fmt.Sprintf(`request_failed{path=%q}`, ctx.Request.URL.Path)
			metrics.GetOrCreateCounter(s).Inc()
		} else {
			s = fmt.Sprintf(`request_success{path=%q}`, ctx.Request.URL.Path)
			metrics.GetOrCreateCounter(s).Inc()
		}
	}
}
```
The middleware function collects total requests count, successful responses count and failed responses count
the failed ones should be alerted

**Important to mention that total request count should be alerted as well in case of significant increase ot detect any kind of DDoS attack or simply high load in time and prevent service from being down**

### For consumers

The metrics are collected on the stage of unmarshalling, in case it is failed the message does not correspond to the
expected structure which means we have troubles with the delivery, this is crucial for the correct functioning of the
services and should be alerted

Code snippet collecting those

```
var message rabbitmq.EmailMessage
err = json.Unmarshal(d.Body, &message)
if err != nil {
	metrics.GetOrCreateCounter(messageConsumedFailMetric).Inc()
	c.logger.Warnf("failed to unmarshal message: %v", err)
	d.Nack(false, false)
	continue
}
metrics.GetOrCreateCounter(fmt.Sprintf(`%v{type=%q}`, messageConsumedSuccessMetric, message.EventType)).Inc()
```

## Process metrics provided by default via victoria metrics package that should be alerted

### Memory metrics
```
go_memstats_alloc_bytes
go_memstats_heap_alloc_bytes
go_memstats_heap_objects
```
This metrics should be alerted as being increased as they 
 - Indicate increase in allocated bytes
 - Indicate high heap memory usage
 - Indicate increase in the number of heap objects

### GC metrics

```
go_gc_cpu_seconds_total
```
This metric should be alerted in case the time cpu spends on garbage collection significantly increasing

### Goroutines metrics

```
go_goroutines
go_threads
```
Those should be alerted in case of significant increase as the high amount of goroutines or threads may be caused by a potential memory leaks


