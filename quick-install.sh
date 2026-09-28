#!/bin/bash
# mmwX Pro - sing-box 订阅管理系统 一键安装命令（简化版）
# 适用于 Debian/Ubuntu Linux 系统

set -e

VERSION="v0.6.0-pro"
GITHUB_REPO="byby5555/mmwx-pro"
VERSION_FILE=".version"
PORT_FILE=".port"

# 检测系统架构
ARCH=$(uname -m)
case "$ARCH" in
    x86_64|amd64)
        BINARY_NAME="mmwx-linux-amd64"
        ;;
    aarch64|arm64)
        BINARY_NAME="mmwx-linux-arm64"
        ;;
    *)
        echo "不支持的系统架构: $ARCH"
        echo "支持的架构: x86_64 (amd64), aarch64 (arm64)"
        exit 1
        ;;
esac

DOWNLOAD_URL="https://github.com/${GITHUB_REPO}/releases/download/${VERSION}/${BINARY_NAME}"

# 安装函数
install() {
    echo "正在下载并安装 mmwX Pro $VERSION ($ARCH)..."

    # 下载
    wget -q --show-progress "$DOWNLOAD_URL" -O mmwx

    # 赋予执行权限
    chmod +x mmwx

    # 创建数据目录
    mkdir -p data

    # 保存版本信息
    echo "$VERSION" > "$VERSION_FILE"

    # 询问端口号（支持非交互式环境）
    echo ""
    if [ -t 0 ]; then
        # 交互式环境
        read -p "请输入端口号（默认 8080，直接回车使用默认值）: " PORT_INPUT
        if [ -z "$PORT_INPUT" ]; then
            PORT=8080
        else
            PORT=$PORT_INPUT
        fi
    else
        # 非交互式环境，使用环境变量或默认值
        PORT=${PORT:-8080}
        echo "使用端口: $PORT"
    fi

    # 保存端口配置
    echo "$PORT" > "$PORT_FILE"

    # 设置环境变量并运行
    export PORT=$PORT
    nohup ./mmwx > mmwx.log 2>&1 &

    # 显示完成信息
    echo ""
    echo "安装完成！"
    echo ""
    echo "访问地址: http://localhost:$PORT"
    echo ""
    echo "更新版本:"
    echo "  curl -sL https://raw.githubusercontent.com/${GITHUB_REPO}/main/quick-install.sh | bash -s update"
    echo ""
    echo "卸载:"
    echo "  curl -sL https://raw.githubusercontent.com/${GITHUB_REPO}/main/quick-install.sh | bash -s uninstall"
    echo ""
}

# 更新函数
update() {
    echo "正在更新 mmwX Pro ($ARCH)..."

    # 检查是否已安装
    if [ ! -f "mmwx" ]; then
        echo "未检测到已安装的 mmwx，请先执行安装"
        exit 1
    fi

    # 显示当前版本
    if [ -f "$VERSION_FILE" ]; then
        CURRENT_VERSION=$(cat "$VERSION_FILE")
        echo "当前版本: $CURRENT_VERSION"
    fi
    echo "目标版本: $VERSION ($ARCH)"
    echo ""

    # 查找并停止运行中的进程
    if pgrep -f "./mmwx" > /dev/null; then
        echo "停止运行中的服务..."
        pkill -f "./mmwx" || true
        sleep 2
    fi

    # 备份当前版本
    if [ -f "mmwx" ]; then
        echo "备份当前版本..."
        cp mmwx mmwx.bak
    fi

    # 下载新版本
    echo "下载新版本..."
    wget -q --show-progress "$DOWNLOAD_URL" -O mmwx

    # 赋予执行权限
    chmod +x mmwx

    # 保存版本信息
    echo "$VERSION" > "$VERSION_FILE"

    # 询问端口号（支持非交互式环境）
    echo ""
    # 尝试读取之前保存的端口号
    SAVED_PORT=""
    if [ -f "$PORT_FILE" ]; then
        SAVED_PORT=$(cat "$PORT_FILE")
    fi

    if [ -t 0 ]; then
        # 交互式环境
        if [ -n "$SAVED_PORT" ]; then
            read -p "请输入端口号（默认 $SAVED_PORT，直接回车使用默认值）: " PORT_INPUT
            if [ -z "$PORT_INPUT" ]; then
                PORT=$SAVED_PORT
            else
                PORT=$PORT_INPUT
            fi
        else
            read -p "请输入端口号（默认 8080，直接回车使用默认值）: " PORT_INPUT
            if [ -z "$PORT_INPUT" ]; then
                PORT=8080
            else
                PORT=$PORT_INPUT
            fi
        fi
    else
        # 非交互式环境，使用环境变量或默认值
        PORT=${PORT:-${SAVED_PORT:-8080}}
        echo "使用端口: $PORT"
    fi

    # 保存端口配置
    echo "$PORT" > "$PORT_FILE"

    # 设置环境变量并运行
    export PORT=$PORT
    nohup ./mmwx > mmwx.log 2>&1 &

    echo ""
    echo "✅ 更新完成！"
    echo ""
    echo "📦 版本: $VERSION"
    echo "🌐 访问地址: http://localhost:$PORT"
    echo ""
    echo "运行服务:"
    echo "  PORT=$PORT ./mmwx"
    echo ""
    echo "后台运行:"
    echo "  PORT=$PORT nohup ./mmwx > mmwx.log 2>&1 &"
    echo ""
    echo "如遇问题可回滚到备份版本:"
    echo "  mv mmwx.bak mmwx"
    echo ""
}

# 卸载函数
uninstall() {
    echo "正在卸载 mmwX Pro..."

    # 检查是否已安装
    if [ ! -f "mmwx" ]; then
        echo "未检测到已安装的 mmwx"
        exit 1
    fi

    # 显示当前版本
    if [ -f "$VERSION_FILE" ]; then
        CURRENT_VERSION=$(cat "$VERSION_FILE")
        echo "当前版本: $CURRENT_VERSION"
        echo ""
    fi

    # 查找并停止运行中的进程
    if pgrep -f "./mmwx" > /dev/null; then
        echo "停止运行中的服务..."
        pkill -f "./mmwx" || true
        sleep 2
        echo "✅ 服务已停止"
        echo ""
    fi

    # 询问是否保留配置和数据
    KEEP_DATA=false
    if [ -t 0 ]; then
        # 交互式环境
        echo "是否保留配置和数据？"
        echo "  1) 完全删除（删除所有文件和数据）"
        echo "  2) 保留数据（保留 data 目录和订阅文件）"
        read -p "请选择 (1/2，默认 2): " CHOICE

        if [ "$CHOICE" = "1" ]; then
            KEEP_DATA=false
        else
            KEEP_DATA=true
        fi
    else
        # 非交互式环境，检查环境变量
        if [ "$KEEP_DATA" != "false" ]; then
            KEEP_DATA=true
        fi
        if [ "$KEEP_DATA" = "true" ]; then
            echo "保留数据模式"
        else
            echo "完全删除模式"
        fi
    fi
    echo ""

    # 删除主程序和版本文件
    echo "删除程序文件..."
    rm -f mmwx mmwx.bak "$VERSION_FILE" "$PORT_FILE" mmwx.log
    echo "✅ 程序文件已删除"
    echo ""

    # 根据选择删除或保留数据
    if [ "$KEEP_DATA" = "false" ]; then
        echo "删除数据和配置..."
        rm -rf data/ subscribes/
        echo "✅ 数据和配置已删除"
        echo ""
        echo "✅ 卸载完成！所有文件已删除"
    else
        echo "保留数据目录: data/"
        echo "保留订阅目录: subscribes/"
        echo ""
        echo "✅ 卸载完成！配置和数据已保留"
        echo ""
        echo "如需重新安装:"
        echo "  curl -sL https://raw.githubusercontent.com/${GITHUB_REPO}/main/quick-install.sh | bash"
    fi
    echo ""
}

# 主函数
main() {
    if [ "$1" = "update" ]; then
        update
    elif [ "$1" = "uninstall" ]; then
        uninstall
    else
        install
    fi
}

# 运行主函数
main "$@"