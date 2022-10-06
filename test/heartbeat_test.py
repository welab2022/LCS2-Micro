import requests
import pytest
import csv


HEARTBEAT_URL = "http://localhost:8080/heartbeat"

test_data_heartbeat = [
    ("200", "Health OK"),
]

@pytest.mark.parametrize("expect_status,  expect_title", test_data_heartbeat)
def test_get_heartbeat_check_status_code_200(expect_status, expect_title):
    response = requests.get(HEARTBEAT_URL)
    response_body = response.json()

    assert response.status_code == 200
    assert response_body["status"] == expect_status 
    assert response_body["title"] == expect_title
    
def test_get_heartbeat_check_status_code_404():
    response = requests.get(HEARTBEAT_URL + "_wrongurl")
    assert response.status_code == 404 
    
# entry function    
if __name__ == "__main__":
    pass
