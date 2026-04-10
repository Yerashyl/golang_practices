import requests
import json
import time

BASE_URL = "http://localhost:8080"

def test_api():
    print("Testing Registration...")
    reg_data = {"email": "test_user@example.com", "password": "password123"}
    r = requests.post(f"{BASE_URL}/register", json=reg_data)
    reg_json = r.json()
    print(f"Status: {r.status_code}, Body: {json.dumps(reg_json)}")
    user_id = reg_json.get("user", {}).get("id")
    print(f"Registered User ID: {user_id}")
    
    print("\nTesting Login...")
    r = requests.post(f"{BASE_URL}/login", json=reg_data)
    print(f"Status: {r.status_code}")
    login_res = r.json()
    access_token = login_res.get("access_token")
    refresh_token = login_res.get("refresh_token")
    print(f"Access Token exists: {access_token is not None}")
    
    headers = {"Authorization": f"Bearer {access_token}"}
    
    print("\nTesting GetMe (should fail - not verified)...")
    r = requests.get(f"{BASE_URL}/users/me", headers=headers)
    print(f"Status: {r.status_code}, Body: {r.text}")

    print("\nTesting Email Verification...")
    r = requests.post(f"{BASE_URL}/verify", json={"email": "test_user@example.com", "code": "1234"})
    print(f"Status: {r.status_code}, Body: {r.text}")

    print("\nTesting GetMe (should succeed now)...")
    r = requests.get(f"{BASE_URL}/users/me", headers=headers)
    print(f"Status: {r.status_code}, Body: {r.text}")
    
    print("\nTesting Promotion (without admin)...")
    r = requests.patch(f"{BASE_URL}/users/promote/{user_id}", headers=headers)
    print(f"Status: {r.status_code}, Body: {r.text}")
    
    print("\nTesting Login as Admin...")
    admin_data = {"email": "admin@example.com", "password": "admin123"}
    r = requests.post(f"{BASE_URL}/login", json=admin_data)
    admin_token = r.json().get("access_token")
    admin_headers = {"Authorization": f"Bearer {admin_token}"}
    
    print(f"\nPromoting User {user_id} as Admin...")
    r = requests.patch(f"{BASE_URL}/users/promote/{user_id}", headers=admin_headers)
    print(f"Status: {r.status_code}, Body: {r.text}")
    
    print("\nTesting Rate Limit on Login (global, limit=10)...")
    limit_reached = False
    for i in range(15):
        r = requests.post(f"{BASE_URL}/login", json=admin_data)
        if r.status_code == 429:
            print(f"Rate limit reached at request {i+1} on /login!")
            limit_reached = True
            break
    if not limit_reached:
        print("Rate limit NOT reached (maybe already used requests).")

    print("\nTesting Refresh Token...")
    r = requests.post(f"{BASE_URL}/users/refresh", json={"refresh_token": refresh_token})
    print(f"Status: {r.status_code}, Body: {r.text}")

if __name__ == "__main__":
    test_api()
