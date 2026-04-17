#!/bin/bash

# 启动脚本

# 设置环境变量
export API_PORT=8080
export DATABASE_PATH=./data/media.db
export TMDB_API_KEY=${TMDB_API_KEY:-""}
export LOG_LEVEL=info

# 确保数据目录存在
mkdir -p ./data ./data/hls ./data/posters ./data/backdrops ./logs

# 启动服务
echo "Starting media server..."
./media-server
