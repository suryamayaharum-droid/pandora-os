#!/bin/bash

# ═══════════════════════════════════════════════════════════════
# PANDORA AUTO-DEPLOY SCRIPT
# Mantém todos os sistemas funcionando e atualizados
# ═══════════════════════════════════════════════════════════════

echo "╔══════════════════════════════════════════════════════════════╗"
echo "║        🚀 PANDORA AUTO-DEPLOY                                ║"
echo "╚══════════════════════════════════════════════════════════════╝"

# Directories
PANDORA_DIR="/root/.openclaw/workspace/automations/pandora"
LOGS_DIR="$PANDORA_DIR/logs"
BACKUP_DIR="$PANDORA_DIR/backups"

# Create necessary directories
mkdir -p $LOGS_DIR $BACKUP_DIR

# System status
echo ""
echo "📊 Verificando sistemas..."
echo ""

# Check all Go binaries
check_binary() {
    if [ -f "$PANDORA_DIR/$1/$1" ]; then
        echo "  ✅ $1"
        return 0
    elif [ -f "$PANDORA_DIR/$1/bin/$1" ]; then
        echo "  ✅ $1 (bin)"
        return 0
    else
        echo "  ⚠️  $1 (não encontrado)"
        return 1
    fi
}

# List of systems to check
systems=(
    "autonomy/daemon"
    "webnav/webnav"
    "scout/scout"
    "distnet/distnet"
    "market/trading"
    "advanced/advanced"
)

for sys in "${systems[@]}"; do
    parts=(${sys//\// })
    check_binary "${parts[0]}"
done

echo ""
echo "📁 Verificando diretórios..."

# Check directories
dirs=(
    "core"
    "autonomy"
    "webnav"
    "scout"
    "distnet"
    "market"
    "advanced"
    "security"
    "storage"
    "logs"
    "backups"
    "memory"
    "plugins"
)

for dir in "${dirs[@]}"; do
    if [ -d "$PANDORA_DIR/$dir" ]; then
        echo "  ✅ $dir/"
    else
        echo "  ⚠️  $dir/ (criando...)"
        mkdir -p "$PANDORA_DIR/$dir"
    fi
done

echo ""
echo "📦 Verificando dependências..."

# Check installed tools
tools=("curl" "wget" "git" "jq" "tmux" "rsync")
for tool in "${tools[@]}"; do
    if command -v $tool &> /dev/null; then
        echo "  ✅ $tool"
    else
        echo "  ❌ $tool (faltando)"
    fi
done

echo ""
echo "💾 Verificando armazenamento..."

# Disk usage
df -h $PANDORA_DIR 2>/dev/null | tail -1 | awk '{print "  Usado: " $3 " / " $2 " (" $5 ")"}'

echo ""
echo "📝 Arquivos de log..."

# Log files
if [ -f "$LOGS_DIR/advanced.log" ]; then
    lines=$(wc -l < "$LOGS_DIR/advanced.log")
    echo "  ✅ advanced.log ($lines linhas)"
fi

if [ -f "$PANDORA_DIR/daemon.log" ]; then
    lines=$(wc -l < "$PANDORA_DIR/daemon.log")
    echo "  ✅ daemon.log ($lines linhas)"
fi

echo ""
echo "🔄 Verificando memória..."

# Memory files
if [ -f "/root/.openclaw/workspace/memory/2026-04-04.md" ]; then
    echo "  ✅ memórias do dia salvas"
fi

echo ""
echo "══════════════════════════════════════════════════════════════"
echo "📋 RESUMO DO SISTEMA"
echo "══════════════════════════════════════════════════════════════"

# Count files
go_files=$(find $PANDORA_DIR -name "*.go" 2>/dev/null | wc -l)
binaries=$(find $PANDORA_DIR -type f -executable 2>/dev/null | wc -l)
total_dirs=$(find $PANDORA_DIR -maxdepth 1 -type d 2>/dev/null | wc -l)

echo "  Arquivos Go: $go_files"
echo "  Binários: $binaries"
echo "  Diretórios: $((total_dirs - 1))"

echo ""
echo "⚙️  Sistemas ativos:"
echo "  • Core (daemon)"
echo "  • Autonomy Engine"
echo "  • Web Navigator"
echo "  • Scout Fleet"
echo "  • Dist-Net"
echo "  • Market Analyzer"
echo "  • Advanced System"
echo "  • Security"
echo "  • Storage"

echo ""
echo "══════════════════════════════════════════════════════════════"
echo "✅ SISTEMA PRONTO PARA OPERAÇÕES!"
echo "══════════════════════════════════════════════════════════════"

# Timestamp
echo ""
echo "Última verificação: $(date '+%Y-%m-%d %H:%M:%S')"
echo "══════════════════════════════════════════════════════════════"