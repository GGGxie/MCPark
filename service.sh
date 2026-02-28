#!/bin/bash

# MCPark 服务管理脚本
# 合并 start.sh + stop.sh，参考 SkillsHub service.sh

set -e

RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
CYAN='\033[0;36m'
NC='\033[0m'

PROJECT_DIR="$(cd "$(dirname "$0")" && pwd)"
BACKEND_DIR="$PROJECT_DIR/backend"
FRONTEND_DIR="$PROJECT_DIR/frontend"

BACKEND_BIN="$BACKEND_DIR/mcpark"
BACKEND_PID_FILE="$BACKEND_DIR/.backend.pid"
BACKEND_LOG="$BACKEND_DIR/backend.log"
FRONTEND_DIST="$FRONTEND_DIR/dist"
STATIC_DIR="$BACKEND_DIR/static"

BACKEND_PORT=8086

# ── 工具函数 ─────────────────────────────────────────────────

log_info() { echo -e "${BLUE}[INFO]${NC} $*"; }
log_ok()   { echo -e "${GREEN}[OK]${NC}   $*"; }
log_warn() { echo -e "${YELLOW}[WARN]${NC} $*"; }
log_err()  { echo -e "${RED}[ERR]${NC}  $*"; }

is_running() {
    local pid_file=$1
    if [ -f "$pid_file" ]; then
        local pid
        pid=$(cat "$pid_file")
        if ps -p "$pid" > /dev/null 2>&1; then
            return 0
        else
            rm -f "$pid_file"
        fi
    fi
    return 1
}

check_port() {
    lsof -i ":$1" > /dev/null 2>&1
}

cleanup_port() {
    local port=$1
    if check_port "$port"; then
        log_warn "端口 $port 已被占用，尝试释放..."
        local pids
        pids=$(lsof -ti ":$port" 2>/dev/null || true)
        [ -n "$pids" ] && kill "$pids" 2>/dev/null || true
        sleep 2
        if check_port "$port"; then
            kill -9 "$pids" 2>/dev/null || true
            sleep 1
        fi
        if check_port "$port"; then
            log_err "无法释放端口 $port"
            return 1
        fi
        log_ok "端口 $port 已释放"
    fi
}

check_deps() {
    command -v go &>/dev/null || { log_err "Go 未安装"; exit 1; }
    command -v npm &>/dev/null || { log_err "npm 未安装"; exit 1; }
}

# ── 帮助 ─────────────────────────────────────────────────────

show_help() {
    echo -e "${BLUE}╔══════════════════════════════════════════════╗${NC}"
    echo -e "${BLUE}║          MCPark 服务管理脚本                 ║${NC}"
    echo -e "${BLUE}╚══════════════════════════════════════════════╝${NC}"
    echo ""
    echo -e "${CYAN}服务控制:${NC}"
    echo "  ./service.sh start           启动所有服务 (后端 + 构建前端)"
    echo "  ./service.sh start backend   仅启动后端"
    echo "  ./service.sh stop            停止所有服务"
    echo "  ./service.sh stop backend    仅停止后端"
    echo "  ./service.sh restart         重启所有服务"
    echo "  ./service.sh restart backend 仅重启后端"
    echo "  ./service.sh status          查看服务状态"
    echo ""
    echo -e "${CYAN}构建:${NC}"
    echo "  ./service.sh build           构建所有 (后端二进制 + 前端静态)"
    echo "  ./service.sh build backend   仅构建后端二进制"
    echo "  ./service.sh build frontend  仅构建前端静态资源"
    echo "  ./service.sh deploy          构建 + 启动 (完整部署)"
    echo ""
    echo -e "${CYAN}开发:${NC}"
    echo "  ./service.sh dev             启动后端开发模式 (go run)"
    echo ""
    echo -e "${CYAN}运维:${NC}"
    echo "  ./service.sh logs            实时查看后端日志"
    echo "  ./service.sh install         安装/更新依赖"
    echo "  ./service.sh clean           清理构建缓存"
    echo "  ./service.sh help            显示帮助"
    echo ""
    echo -e "${CYAN}端口信息:${NC}"
    echo "  后端 API: :$BACKEND_PORT"
    echo "  访问地址: http://localhost:$BACKEND_PORT"
    echo ""
}

# ── 后端 ─────────────────────────────────────────────────────

build_backend() {
    log_info "构建后端二进制..."
    cd "$BACKEND_DIR"
    go mod tidy
    go build -o mcpark ./cmd/main.go
    log_ok "后端构建成功 → $BACKEND_BIN"
}

start_backend() {
    log_info "启动后端服务..."
    if is_running "$BACKEND_PID_FILE"; then
        log_warn "后端已在运行 (PID: $(cat "$BACKEND_PID_FILE"))"
        return 0
    fi
    cleanup_port "$BACKEND_PORT"

    [ -f "$BACKEND_BIN" ] || build_backend

    cd "$BACKEND_DIR"
    mkdir -p data
    nohup "$BACKEND_BIN" > "$BACKEND_LOG" 2>&1 &
    echo $! > "$BACKEND_PID_FILE"
    sleep 2

    if is_running "$BACKEND_PID_FILE"; then
        log_ok "后端启动成功  http://localhost:$BACKEND_PORT  (PID: $(cat "$BACKEND_PID_FILE"))"
    else
        log_err "后端启动失败，查看日志: $BACKEND_LOG"
        tail -20 "$BACKEND_LOG"
        exit 1
    fi
}

stop_backend() {
    log_info "停止后端服务..."
    if is_running "$BACKEND_PID_FILE"; then
        local pid
        pid=$(cat "$BACKEND_PID_FILE")
        kill "$pid" 2>/dev/null || true
        sleep 2
        ps -p "$pid" &>/dev/null && kill -9 "$pid" 2>/dev/null || true
        rm -f "$BACKEND_PID_FILE"
        log_ok "后端已停止"
    else
        # 兜底：杀死所有 mcpark 进程
        pkill -f "mcpark" 2>/dev/null || true
        lsof -ti ":$BACKEND_PORT" 2>/dev/null | xargs kill -9 2>/dev/null || true
        log_warn "后端未运行（已尝试清理残留进程）"
    fi
}

dev_backend() {
    log_info "启动后端开发模式 (go run)..."
    cd "$BACKEND_DIR"
    mkdir -p data
    go run ./cmd/main.go
}

# ── 前端 ─────────────────────────────────────────────────────

build_frontend() {
    log_info "构建前端静态资源..."
    cd "$FRONTEND_DIR"
    [ ! -d "node_modules" ] && npm install --legacy-peer-deps
    npm run build
    if [ -d "$FRONTEND_DIST" ]; then
        # 复制到后端 static 目录
        mkdir -p "$STATIC_DIR"
        rm -rf "${STATIC_DIR:?}"/*
        cp -r "$FRONTEND_DIST"/. "$STATIC_DIR"/
        log_ok "前端构建成功 → $STATIC_DIR ($(du -sh "$STATIC_DIR" | cut -f1))"
    else
        log_err "前端构建失败"
        exit 1
    fi
}

# ── 聚合操作 ─────────────────────────────────────────────────

start_all() {
    check_deps
    build_frontend
    start_backend
    echo ""
    show_status
}

stop_all() {
    log_info "停止所有服务..."
    stop_backend
}

restart_all() {
    stop_all
    sleep 2
    start_all
}

deploy() {
    log_info "完整部署流程: 构建 → 启动..."
    check_deps
    build_backend
    build_frontend
    stop_all
    sleep 1
    start_backend
    echo ""
    show_status
    log_ok "部署完成 → http://localhost:$BACKEND_PORT"
}

install_deps() {
    log_info "安装依赖..."
    cd "$BACKEND_DIR" && go mod tidy
    log_ok "Go 依赖更新完成"
    cd "$FRONTEND_DIR" && npm install
    log_ok "Node 依赖安装完成"
}

clean_all() {
    log_info "清理构建缓存..."
    rm -f "$BACKEND_BIN"
    rm -rf "$FRONTEND_DIST" "$FRONTEND_DIR/node_modules/.cache"
    rm -rf "$STATIC_DIR"
    log_ok "清理完成"
}

show_status() {
    echo -e "${BLUE}=== MCPark 服务状态 ===${NC}"
    if is_running "$BACKEND_PID_FILE"; then
        echo -e "  ${GREEN}●${NC} 后端  运行中  http://localhost:$BACKEND_PORT  (PID: $(cat "$BACKEND_PID_FILE"))"
    else
        echo -e "  ${RED}○${NC} 后端  未运行"
    fi

    echo ""
    echo -e "${CYAN}构建产物:${NC}"
    [ -f "$BACKEND_BIN" ] \
        && echo -e "  后端二进制: $(du -sh "$BACKEND_BIN" | cut -f1)" \
        || echo -e "  后端二进制: ${YELLOW}未构建${NC}"
    [ -d "$STATIC_DIR" ] \
        && echo -e "  前端静态:   $(du -sh "$STATIC_DIR" | cut -f1)" \
        || echo -e "  前端静态:   ${YELLOW}未构建${NC}"

    echo ""
    echo -e "${CYAN}访问地址:${NC}"
    echo "  http://localhost:$BACKEND_PORT"
    echo ""
}

show_logs() {
    if [ -f "$BACKEND_LOG" ]; then
        log_info "实时查看后端日志 (Ctrl+C 退出)..."
        tail -f "$BACKEND_LOG"
    else
        log_warn "暂无日志文件: $BACKEND_LOG"
    fi
}

# ── 入口 ─────────────────────────────────────────────────────

main() {
    case "${1:-help}" in
        start)
            case "${2:-all}" in
                all|"") start_all ;;
                backend) start_backend ;;
                *) log_err "未知参数: $2"; show_help ;;
            esac ;;
        stop)
            case "${2:-all}" in
                all|"") stop_all ;;
                backend) stop_backend ;;
                *) log_err "未知参数: $2"; show_help ;;
            esac ;;
        restart)
            case "${2:-all}" in
                all|"") restart_all ;;
                backend) stop_backend; sleep 1; start_backend ;;
                *) log_err "未知参数: $2"; show_help ;;
            esac ;;
        build)
            case "${2:-all}" in
                all|"")  build_backend; build_frontend ;;
                backend) build_backend ;;
                frontend) build_frontend ;;
                *) log_err "未知参数: $2"; show_help ;;
            esac ;;
        dev)  dev_backend ;;
        deploy) deploy ;;
        status) show_status ;;
        logs)   show_logs ;;
        install) install_deps ;;
        clean) clean_all ;;
        help|-h|--help) show_help ;;
        *) log_err "未知命令: $1"; echo ""; show_help; exit 1 ;;
    esac
}

main "$@"
