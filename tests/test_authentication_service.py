import requests
import pytest
import csv

AUTH_SERVICE = "http://localhost:8081/api/auth"


def test_get_heartbeat_check_status_code_200():
    response = requests.get(AUTH_SERVICE + "/heartbeat")
    assert response.status_code == 200

def test_get_heartbeat_check_status_code_404():
    response = requests.get(AUTH_SERVICE + "_wrongurl")
    assert response.status_code == 404 

# response_body = response.json()