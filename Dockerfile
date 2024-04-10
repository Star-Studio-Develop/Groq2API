# 使用官方 Python 运行时作为父镜像
FROM python:3.9-slim

# 将工作目录设置为 /app
WORKDIR /app

# 将当前目录内容复制到位于 /app 中的容器里
COPY . /app

# 安装 requirements.txt 中指定的任何所需包
RUN pip install --no-cache-dir -r requirements.txt

# 让世界知道在容器内部应用运行在哪个端口
EXPOSE 8000

# 定义环境变量
ENV NAME World

# 在容器启动时运行的命令
CMD ["uvicorn", "main:app", "--host", "0.0.0.0", "--port", "8000"]
