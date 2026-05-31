import os
import time
import json
import logging
from fastapi import FastAPI, Request
from pydantic import BaseModel
from crawl4ai import AsyncWebCrawler, BrowserConfig, CrawlerRunConfig, CacheMode

from opentelemetry import trace
from opentelemetry.sdk.resources import Resource
from opentelemetry.sdk.trace import TracerProvider
from opentelemetry.exporter.otlp.proto.http.trace_exporter import OTLPSpanExporter
from opentelemetry.sdk.trace.export import BatchSpanProcessor
from opentelemetry.instrumentation.fastapi import FastAPIInstrumentor

app = FastAPI(title="Searchbase Crawl4AI Internal Worker")
log_levels = {
    "debug": logging.DEBUG,
    "info": logging.INFO,
    "warn": logging.WARNING,
    "error": logging.ERROR,
}


OTEL_TRACING = os.environ.get("SEARCHBASE_ENABLE_TRACING", "false").lower()
OTEL_ENDPOINT = os.environ.get("SEARCHBASE_OTLP_ENDPOINT")
LOG_LEVEL = os.environ.get("CRAWLER_LOG_LEVEL", "info").lower()
LOG_FORMAT = os.environ.get("CRAWLER_LOG_FORMAT", "json").lower()


class JSONFOrmatter(logging.Formatter):
    """Format log records as compact JSON for structured logging backends."""

    def format(self, record):
        log_record = {
            "time": self.formatTime(record, self.datefmt),
            "level": record.levelname,
            "msg": record.getMessage(),
        }

        if hasattr(record, "attrs"):
            log_record.update(record.attrs)
        return json.dumps(log_record)


class TextFormatter(logging.Formatter):
    """Format log records as key-value text for local debugging."""

    def format(self, record):
        log_str = f'time={self.formatTime(record, self.datefmt)} level={record.levelname} msg="{record.getMessage()}"'
        if hasattr(record, "attrs"):
            for k, v in record.attrs.items():
                log_str += f" {k}={v}"
            return log_str


class ExtractRequest(BaseModel):
    """Internal request body for extracting one URL."""

    url: str
    js_render: bool = False


class ExtractResponse(BaseModel):
    """Internal response body returned after extraction."""

    markdown: str
    success: bool
    error: str = ""


## Logger setup
logger = logging.getLogger("crawl-worker")
logger.setLevel(log_levels.get(LOG_LEVEL, logging.INFO))
handler = logging.StreamHandler()
handler.setFormatter(JSONFOrmatter() if LOG_FORMAT == "json" else TextFormatter())
logger.addHandler(handler)
logging.getLogger("uvicorn.access").disabled = True

## OTEL Tracing setup
if OTEL_TRACING == "true":
    resource = Resource(attributes={"service.name": "crawl-worker"})
    provider = TracerProvider(resource=resource)
    exporter = OTLPSpanExporter(endpoint=OTEL_ENDPOINT)
    processor = BatchSpanProcessor(exporter)
    provider.add_span_processor(processor)
    trace.set_tracer_provider(provider)
    FastAPIInstrumentor.instrument_app(app)


@app.middleware("http")
async def slogger_middleware(request: Request, call_next):
    """Log request metadata without raw IPs or user-identifying headers."""

    start = time.time()
    status_code = 500

    try:
        response = await call_next(request)
        status_code = response.status_code
        return response
    finally:
        latency = time.time() - start

        attrs = {
            "method": request.method,
            "path": request.url.path,
            "status": status_code,
            "latency": f"{latency * 1000:.2f}ms",
        }

        current_span = trace.get_current_span()
        if current_span and current_span.get_span_context().is_valid:
            span_context = current_span.get_span_context()
            attrs["trace_id"] = f"{span_context.trace_id:032x}"
            attrs["span_id"] = f"{span_context.span_id:016x}"

        # Trust the edge proxy for region data, never log the raw IP
        country = request.headers.get("CF-IPCountry")
        if country:
            attrs["country"] = country

        msg = "request handled"
        level = logging.INFO

        if status_code >= 500:
            level = logging.ERROR
            msg = "server error"

        logger.log(level, msg, extra={"attrs": attrs})


@app.post("/extract", response_model=ExtractResponse)
async def extract_content(req: ExtractRequest):
    """Fetch a URL with crawl4ai and return LLM-friendly Markdown."""

    span = trace.get_current_span()
    try:
        browser_config = BrowserConfig(verbose=False)
        run_config = CrawlerRunConfig(
            verbose=False,
            log_console=False,
            cache_mode=CacheMode.BYPASS,
            js_code="" if not req.js_render else None,
        )
        async with AsyncWebCrawler(config=browser_config) as crawler:
            result = await crawler.arun(url=req.url, config=run_config)

            if result.success:
                return ExtractResponse(markdown=result.markdown, success=True)
            else:
                span.set_status(trace.StatusCode.ERROR, "crawling failed")
                span.record_exception(Exception(result.error_message))
                return ExtractResponse(
                    markdown="", success=False, error=result.error_message
                )

    except Exception as e:
        span.set_status(trace.StatusCode.ERROR, "unhandled exception during crawl")
        span.record_exception(e)
        return ExtractResponse(markdown="", success=False, error=str(e))


if __name__ == "__main__":
    import uvicorn

    uvicorn.run(app, host="0.0.0.0", port=8000)
