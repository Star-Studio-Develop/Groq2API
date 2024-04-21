# Groq2API

## Installation

```bash

docker run -d -p 8080:8080 ghcr.io/star-studio-develop/groq2api:latest

```

### Vercel部署
[![Deploy with Vercel](https://vercel.com/button)](https://vercel.com/new/clone?repository-url=https%3A%2F%2Fgithub.com%2FStar-Studio-Develop%2FGroq2API&project-name=Groq2API&repository-name=Groq2API)

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
  - llama3-8b-8192
  - llama3-70b-8192
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
  --header 'Authorization: Bearer stytch_session' \
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

### stytch_session 获取方法
![image|690x233](https://cdn.linux.do/uploads/default/original/3X/c/c/cc5bb06024b2fc93581227e16b5a5e3e220d159c.png)
