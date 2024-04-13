# Groq2API

## Installation

```bash

docker run -d -p 8080:8080 ghcr.io/star-studio-develop/groq2api:latest

```

### Koyeb部署
[![Deploy to Koyeb](https://www.koyeb.com/static/images/deploy/button.svg)](https://app.koyeb.com/deploy?type=docker&name=groq2api&ports=8080;http;/&image=ghcr.io/star-studio-develop/groq2api)

### Render部署
[![Deploy](https://render.com/images/deploy-to-render-button.svg)](https://render.com/deploy)

### Railway部署
[![Deploy on Railway](https://railway.app/button.svg)](https://railway.app/template/kaB56J)

## Usage

可选参数列表

- `model` 模型名称
  - gemma-7b-it
  - mixtral-8x7b-32768
  - llama2-70b-4096
- `stream` 是否流式输出
  - true
  - false
- `max_tokens` 最大生成长度
  - 4096 (llama2-70b-4096) 
  - 8192 (gemma-7b-it)
  - 32768 (mixtral-8x7b-32768)
- `message`
  - `role` 消息角色
    - user
    - assistant
```bash

curl --request POST \
  --url http://127.0.0.1:8080/v1/chat/completions \
  --header 'Authorization: Bearer change-it-to-your-refresh-token' \
  --data '{
  "messages": [
    {
      "role": "user",
      "content": "hi"
    }
  ],
  "model": "mixtral-8x7b-32768",
  "max_tokens": 4096,
  "stream": true
}'

```

![image](https://github.com/Star-Studio-Develop/Groq2API/assets/148524140/adedf992-864a-47b1-9201-d53717befd4a)

