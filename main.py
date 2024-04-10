from fastapi import FastAPI
from pydantic import BaseModel
import httpx
import asyncio
from typing import List, Dict, Union

app = FastAPI()

class Message(BaseModel):
    role: str
    content: str

class CompletionRequest(BaseModel):
    model: str
    messages: List[Message]
    stream: bool

refresh_token = 'cHVibGljLXRva2VuLWxpdmUtMjZhODlmNTktMDlmOC00OGJlLTkxZmYtY2U3MGU2MDAwY2I1OnpUVTZoNnBZcTZhMFRacEJlSnY3R3pwSEs1OEFnZElLdjBqa3hzVnNxVExZ'

@app.post("/v1/chat/completions")
async def fetch_stream(request: CompletionRequest):
    async with httpx.AsyncClient() as client:
        # 第一次请求获取session_jwt
        headers_auth = {
        'accept': '*/*',
        'accept-language': 'zh-CN,zh;q=0.9',
        'authorization': f'Basic {refresh_token}',
        'cache-control': 'no-cache',
        'content-type': 'application/json',
        'origin': 'https://groq.com',
        'pragma': 'no-cache',
        'referer': 'https://groq.com/',
        'user-agent': 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36',
        'x-sdk-client': 'eyJldmVudF9pZCI6ImV2ZW50LWlkLWQ2ZTJmZGM3LTVhYjYtNGZjNy05NjQxLTQ5ZTNjMDZkNDgxYyIsImFwcF9zZXNzaW9uX2lkIjoiYXBwLXNlc3Npb24taWQtNzdkMGRhYjMtZjUwMi00MWQyLWJiMTUtYjMxYzBlOGI5MzY0IiwicGVyc2lzdGVudF9pZCI6InBlcnNpc3RlbnQtaWQtMDVjZDVlZTEtYTVkNC00YjA5LWI3NmEtYmNkNTU4MTY0ODkxIiwiY2xpZW50X3NlbnRfYXQiOiIyMDI0LTA0LTA5VDEzOjEyOjUxLjk5NloiLCJ0aW1lem9uZSI6IkFzaWEvU2hhbmdoYWkiLCJzdHl0Y2hfdXNlcl9pZCI6InVzZXItbGl2ZS1kZmJlODE4My00MWMzLTQyM2EtODBlNy02MWMzNWRkODQ5ODMiLCJzdHl0Y2hfc2Vzc2lvbl9pZCI6InNlc3Npb24tbGl2ZS1iNDhmN2RjNy05MGI4LTRkNDYtYWExYi1jYmVmZDBkN2ZjMTQiLCJhcHAiOnsiaWRlbnRpZmllciI6Imdyb3EuY29tIn0sInNkayI6eyJpZGVudGlmaWVyIjoiU3R5dGNoLmpzIEphdmFzY3JpcHQgU0RLIiwidmVyc2lvbiI6IjQuNS4zIn19',
        'x-sdk-parent-host': 'https://groq.com',
    }
        response = await client.post('https://web.stytch.com/sdk/v1/sessions/authenticate', headers=headers_auth, json={})
        session_jwt = response.json()['data']['session_jwt']

        # 使用session_jwt获取org_id
        headers_org = {
        "accept": "*/*",
        "accept-language": "zh-CN,zh;q=0.9",
        "cache-control": "no-cache",
        "content-type": "application/json",
        "pragma": "no-cache",
        "sec-ch-ua": "\"Google Chrome\";v=\"123\", \"Not:A-Brand\";v=\"8\", \"Chromium\";v=\"123\"",
        "sec-ch-ua-mobile": "?0",
        "sec-ch-ua-platform": "\"macOS\"",
        "sec-fetch-dest": "empty",
        "sec-fetch-mode": "cors",
        "sec-fetch-site": "same-site",
        "authorization": f"Bearer {session_jwt}",
}
        response = await client.get("https://api.groq.com/platform/v1/user/profile", headers=headers_org)
        org_id = response.json()['user']['orgs']['data'][0]['id']

        # 使用org_id和session_jwt进行最终的流式请求
        headers_stream = {
        'accept': 'application/json',
        'accept-language': 'zh-CN,zh;q=0.9',
        'cache-control': 'no-cache',
        'content-type': 'application/json',
        'groq-app': 'chat',
        'groq-organization': org_id,
        'origin': 'https://groq.com',
        'pragma': 'no-cache',
        'referer': 'https://groq.com/',
        'sec-ch-ua': '"Google Chrome";v="123", "Not:A-Brand";v="8", "Chromium";v="123"',
        'sec-ch-ua-mobile': '?0',
        'sec-ch-ua-platform': '"macOS"',
        'sec-fetch-dest': 'empty',
        'sec-fetch-mode': 'cors',
        'sec-fetch-site': 'same-site',
        'user-agent': 'Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36',
        'x-stainless-arch': 'unknown',
        'x-stainless-lang': 'js',
        'x-stainless-os': 'Unknown',
        'x-stainless-package-version': '0.3.2',
        'x-stainless-runtime': 'browser:chrome',
        'x-stainless-runtime-version': '123.0.0',
        'authorization': f"Bearer {session_jwt}"
    }
        messages_serializable = [message.dict() for message in request.messages]

        data = {
    "model": request.model,
    "messages": messages_serializable,  # 使用可序列化的messages代替原来的版本
    "temperature": 0.2,
    "max_tokens": 2048,
    "top_p": 0.8,
    "stream": request.stream
}
        async with client.stream('POST', 'https://api.groq.com/openai/v1/chat/completions', headers=headers_stream, json=data) as response:
            results = []
            async for chunk in response.aiter_text():
                results.append(chunk)
            return {"stream_results": results}