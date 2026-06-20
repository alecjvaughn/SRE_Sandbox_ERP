import http.server
import socketserver
from prometheus_client import start_http_server
import logging

logger = logging.getLogger(__name__)

class HealthCheckHandler(http.server.SimpleHTTPRequestHandler):
    def do_GET(self):
        if self.path == '/healthz':
            self.send_response(200)
            self.send_header('Content-type', 'text/plain')
            self.end_headers()
            self.wfile.write(b"OK")
        else:
            self.send_response(404)
            self.end_headers()
    
    # Override log_message to prevent test output clutter
    def log_message(self, format, *args):
        pass

class ReusableTCPServer(socketserver.TCPServer):
    allow_reuse_address = True

def run_health_server(port=8081):
    handler = HealthCheckHandler
    with ReusableTCPServer(("", port), handler) as httpd:
        logger.info(f"Health server serving at port {port}")
        httpd.serve_forever()

def start_metrics_server(port=8000):
    start_http_server(port)
    logger.info(f"Metrics server started on port {port}")

if __name__ == "__main__":
    import threading
    from consumer import consume_orders
    
    start_metrics_server(8000)
    
    # Start consumer in a background thread
    consumer_thread = threading.Thread(target=consume_orders, daemon=True)
    consumer_thread.start()
    
    run_health_server(8081)
