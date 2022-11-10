import requests
import pytest
import names
import csv

AUTH_SERVICE = "http://localhost:8081/api/auth"
API_KEY = "sWOmNsF8Ht9lE9wMU9cW7w==n" 

def helper_login(credentials):
    url = AUTH_SERVICE + "/signin"
    headers = {"Content-Type": "application/json; charset=utf-8", "X-API-Key": API_KEY}
    data = credentials 
    return requests.post(url, headers=headers, json=data)

def helper_get_cookie(response):
    cookie_session = response.headers['Set-Cookie'].split('=')[0] 
    cookie_token = response.headers['Set-Cookie'].split('=')[1].split(';')[0] 
    # print(response.cookies["lcs2_session_token"])
    cookies_dict = {cookie_session: cookie_token}
    return cookies_dict 

def test_get_heartbeat_check_status_code_200():
    response = requests.get(AUTH_SERVICE + "/heartbeat")
    assert response.status_code == 200

def test_get_wrongurl_heartbeat_check_status_code_404():
    response = requests.get(AUTH_SERVICE + "_wrongurl")
    assert response.status_code == 404 

def test_post_admin_signin_check_status_code_202():
    url = AUTH_SERVICE + "/signin"
   
    headers = {"Content-Type": "application/json; charset=utf-8", "X-API-Key": API_KEY}
    
    data = {
        "email": "admin@example.com", 
        "password": "verysecret", 
    }
    
    response = requests.post(url, headers=headers, json=data)
    assert response.status_code == 202
     
def test_post_admin_logout_check_status_code_200():
    
    # sigin request with an admin credential
    credentials = {
        "email": "admin@example.com", 
        "password": "verysecret", 
    }
    res = helper_login(credentials)
    cookies = helper_get_cookie(res)
    
    # logout request
    url = AUTH_SERVICE + "/logout"
    headers = {"Content-Type": "application/json; charset=utf-8", "X-API-Key": API_KEY}
    data = { "email": "admin@example.com" }
    response = requests.post(url, headers=headers, json=data, cookies=cookies)
    assert response.status_code == 200

def test_get_admin_list_users_check_status_code_200():
    # sigin request with an admin credential
    credentials = {
        "email": "admin@example.com", 
        "password": "verysecret", 
    }
    res = helper_login(credentials)
    cookies = helper_get_cookie(res)
    
    # /listusers request
    url = AUTH_SERVICE + "/listusers"
    headers = {"Content-Type": "application/json; charset=utf-8", "X-API-Key": API_KEY}
    response = requests.get(url, headers=headers, cookies=cookies)
    assert response.status_code == 200

def test_post_admin_avatar_check_status_code_200():
    # sigin request with an admin credential
    credentials = {
        "email": "admin@example.com", 
        "password": "verysecret", 
    }
    res = helper_login(credentials)
    cookies = helper_get_cookie(res)
    
    # /upload post request
    url = AUTH_SERVICE + "/upload"
 
    headers = { "X-API-Key": API_KEY}
    
    image = open('./avatars/logo.png', 'rb')
    files= {'file': image}
    response = requests.post(url, files=files, headers=headers, cookies=cookies)
    assert response.status_code == 200
      
def test_get_admin_avatar_check_status_code_200():
    # sigin request with an admin credential
    credentials = {
        "email": "admin@example.com", 
        "password": "verysecret", 
    }
    res = helper_login(credentials)
    cookies = helper_get_cookie(res)
    
    # get request => /avatar/admin@example.com
    url = AUTH_SERVICE + "/avatar/admin@example.com"
    headers = {"Content-Type": "application/json; charset=utf-8", "X-API-Key": API_KEY}
    response = requests.get(url,  headers=headers, cookies=cookies)
    assert response.status_code == 200
    
def test_post_admin_adduser_check_status_code_200():
    # sigin request with an admin credential
    credentials = {
        "email": "admin@example.com", 
        "password": "verysecret", 
    }
    res = helper_login(credentials)
    cookies = helper_get_cookie(res)
    
    # post request => /adduser
    url = AUTH_SERVICE + "/adduser"
    headers = {"Content-Type": "application/json; charset=utf-8", "X-API-Key": API_KEY}
    
    rand_fname = names.get_first_name()
    rand_lname = names.get_last_name()
    email = "{fname}.{lname}@test.com".format(fname=rand_fname.lower(), lname=rand_lname.lower())
    data = {
        "email": email,
        "first_name": rand_fname,
        "last_name": rand_lname,
        "password": "12345"
    }
    response = requests.post(url,  headers=headers, json=data, cookies=cookies)
    # print(response.content)
    assert response.status_code == 200

def test_post_admin_existed_adduser_check_status_code_500():
    # sigin request with an admin credential
    credentials = {
        "email": "admin@example.com", 
        "password": "verysecret", 
    }
    res = helper_login(credentials)
    cookies = helper_get_cookie(res)
    
    # post request => /adduser
    url = AUTH_SERVICE + "/adduser"
    headers = {"Content-Type": "application/json; charset=utf-8", "X-API-Key": API_KEY}
    data = {
        "email": "admin@example.com", 
        "first_name": "Admin",
        "last_name": "User",
        "password": "12345"
    }
    response = requests.post(url,  headers=headers, json=data, cookies=cookies)
    # print(response.content)
    # print(response.status_code)
    assert response.status_code == 500
    