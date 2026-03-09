# Desafio Clean Architecture 

Projeto de crud de ordens com criação e Listagem de ordens.

## 🚀 Descrição

Projeto desenvolvido em Go que implementa um serviço de ordens com pelos canais (API RESTful, GraphQl e Grpc) para criar e listar.

## 📋 Pré-requisito

- Go 1.19 ou superior
- Git

## Dependências e sugestões
- Docker - Subir a infra
- Evans - Se comunicar com o serviço de grpc

## 🏃 Como Executar 

1. Execute o comando: 

   ```bash
   docker compose up -d 
   ```

O servidor estará disponível em:
-Rest - `http://localhost:8000`
-Graphql - `http://localhost:8080/graphql`
-Grpc - `http://localhost:50051`

## 🧪 Como Rodar os Testes

Para executar todos os testes do projeto:

```bash
go test ./...
```
Para rodar testes com cobertura de código:

```bash
go test -cover ./...
```
### Usando a extensão REST Client

1. **Instale a extensão** no VS Code:
   - Procure por "REST Client" (publicada por Huachao Mao)
   - Ou execute: `ext install humao.rest-client`

2. **Use o arquivo** `api/create_order.http` incluído no projeto:
   - Abra o arquivo `api/create_order.http` ou o `api/get_all.http`
   - Clique em "Send Request" (ou use `Ctrl+Alt+R`)
   - Veja a resposta no painel de output

## 📞 Contato

Desenvolvido por Luana Andrade - luanaands@gmail.com

---

**Aproveite! 🚀**
