# Groq2API

## Installation

```bash

docker run -d -p 8080:8080  ghcr.io/star-studio-develop/groq2api:latest

```

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

