import urllib.request
import urllib.error
import threading
import time
import pytest
from app import run_health_server

@pytest.fixture(scope="module")
def health_server():
    # Start the server in a background thread for testing
    server_thread = threading.Thread(target=run_health_server, daemon=True)
    server_thread.start()
    # Give it a moment to start
    time.sleep(1)
    yield

def test_healthz_endpoint(health_server):
    req = urllib.request.Request('http://localhost:8081/healthz')
    try:
        with urllib.request.urlopen(req, timeout=2) as response:
            assert response.status == 200
            body = response.read().decode('utf-8')
            assert "OK" in body
    except urllib.error.URLError as e:
        pytest.fail(f"Health check failed: {e}")

def test_notfound_endpoint(health_server):
    req = urllib.request.Request('http://localhost:8081/invalid')
    try:
        with urllib.request.urlopen(req, timeout=2) as response:
            pytest.fail("Expected 404 error")
    except urllib.error.HTTPError as e:
        assert e.code == 404
