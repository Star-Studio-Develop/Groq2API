import requests

url = "https://YOUR_URL/v1/chat/completions"
headers = {
    "Authorization": "Bearer your-refresh-token",
    "Content-Type": "application/json"
}
data = {
    "messages": [
        {
            "role": "user",
            "content": "hi"
        }
    ],
    "model": "mixtral-8x7b-32768",
    "max_tokens": 4096,
    "stream": True
}

response = requests.post(url, headers=headers, json=data)
print(response.text)
