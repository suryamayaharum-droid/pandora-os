# Agent MVP com Human-in-the-Loop (HITL)

Implementação mínima de um agente funcional com aprovação humana antes de execução externa.

## Fluxo

1. Recebe objetivo.
2. Gera plano (via endpoint OpenAI-compatible opcional).
3. Solicita aprovação humana.
4. Executa em modo seguro.
5. Salva relatório em `hitl_run_state.json`.

## Execução CLI

```bash
cd agent
go run ./hitl
```

## Acesso por HTML

```bash
cd agent
go run ./hitl -web -addr :8080
```

Depois, abra no navegador:

- http://localhost:8080


## Acesso rápido (1 comando)

```bash
cd agent/hitl
./start_web.sh 8080
```

Esse script já mostra os links para acesso local e pela rede.

## Integração com modelos open-source (opcional)

Defina variáveis para usar servidor compatível com OpenAI (ex.: vLLM, Ollama compat layer):

```bash
export OPENAI_COMPAT_BASE_URL="http://localhost:8000/v1"
export OPENAI_COMPAT_MODEL="Qwen/Qwen2.5-7B-Instruct"
export OPENAI_COMPAT_API_KEY=""
```

Sem essas variáveis, o sistema usa plano de fallback local.
