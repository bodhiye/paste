#!/bin/bash

# 定义颜色
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'

# 显示菜单
show_menu() {
    echo -e "${GREEN}Docker Compose 管理脚本${NC}"
    echo "------------------------"
    echo -e "${BLUE}1. 启动所有服务"
    echo "2. 停止服务"
    echo "3. 停止服务并删除所有数据"
    echo "4. 退出"
}

# 启动服务
start_services() {
    echo -e "${GREEN}正在启动所有服务...${NC}"
    docker-compose up -d
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}服务启动成功！${NC}"
        echo -e "${BLUE}您可以通过以下地址访问服务：${NC}"
        echo -e "${BLUE}➜ http://127.0.0.1:80${NC}"
    else
        echo -e "${RED}服务启动失败！${NC}"
    fi
}

# 停止服务
stop_services() {
    echo -e "${YELLOW}正在停止服务...${NC}"
    docker-compose down
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}服务已停止！${NC}"
    else
        echo -e "${RED}停止服务失败！${NC}"
    fi
}

# 停止服务并删除数据
stop_and_clean() {
    echo -e "${RED}警告：这将删除所有容器和数据！${NC}"
    read -p "确定要继续吗？(y/n): " confirm
    if [ "$confirm" = "y" ] || [ "$confirm" = "Y" ]; then
        echo -e "${YELLOW}正在停止服务并清理数据...${NC}"
        docker-compose down -v
        # 删除MongoDB数据目录
        rm -rf ../mongodb/db/*
        if [ $? -eq 0 ]; then
            echo -e "${GREEN}服务已停止，所有数据已清理！${NC}"
        else
            echo -e "${RED}操作失败！${NC}"
        fi
    else
        echo -e "${YELLOW}操作已取消${NC}"
    fi
}

# 主循环
while true; do
    show_menu
    read -p "请选择操作 (1-4): " choice
    case $choice in
        1)
            start_services
            ;;
        2)
            stop_services
            ;;
        3)
            stop_and_clean
            ;;
        4)
            echo -e "${GREEN}再见！${NC}"
            exit 0
            ;;
        *)
            echo -e "${RED}无效的选择，请重试${NC}"
            ;;
    esac
    echo
done