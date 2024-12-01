# Rate Limiter em Go

## Objetivo
Este projeto implementa um rate limiter em Go que controla o número de requisições por segundo com base em um endereço IP ou token de acesso. Ele foi projetado para ser configurável e escalável, utilizando Redis para persistência e permitindo fácil troca de mecanismo de armazenamento via estratégia de persistência.

## Como Funciona
### Funcionalidade Principal
- **Limitação por IP**: O número de requisições permitidas por segundo para um determinado endereço IP é restrito. Caso o limite seja atingido, o IP será bloqueado até que o tempo de bloqueio expire.  
- **Limitação por Token de Acesso**: Cada token de acesso pode ter um limite de requisições por segundo. As configurações de limite de token têm prioridade sobre o limite de IP, ou seja, se o limite de token for maior, ele será o fator determinante para o bloqueio.

- **Persistência com Redis**: O estado do rate limiter é armazenado no Redis para garantir que as informações sobre os limites e bloqueios sejam persistentes entre reinícios do servidor.

## Configuração: 
O rate limiter pode ser configurado através de variáveis de ambiente ou um arquivo `.env`. 
As configurações incluem o número máximo de requisições por segundo para o IP e o token, o tempo de expiração para bloqueios, e o tipo de mecanismo de persistência a ser utilizado.

### Subir imagem Docker
```batch
docker build -t
```
## Como Executar os Testes
### Requisitos

Para rodar os testes, você precisará dos seguintes pacotes Go instalados:

- `github.com/stretchr/testify`: Usado para asserções e mocks.

Você pode instalar as dependências com o comando:

```bash
go get -u github.com/stretchr/testify
```
# Rodando os Testes:
  Executando os Testes: Após garantir que as dependências estejam instaladas, você pode rodar os testes com o comando:

```batch
  go test -v
```
 O parâmetro -v irá fornecer uma saída detalhada do processo de execução dos testes.
 Estrutura dos Testes: Os testes estão localizados no arquivo limiter_test.go e são baseados no pacote testify/assert para asserções e testify/mock para simular o comportamento do Redis e do rate limiter.

Testes Disponíveis:
    -> Teste de Limitação por IP (TestRateLimiter_LimitByIP): Verifica se as requisições de um IP que excedem o limite configurado são corretamente bloqueadas.
    -> Teste de Limitação por Token (TestRateLimiter_LimitByToken): Verifica se as requisições de um token que excedem o limite configurado são corretamente bloqueadas.
    -> Teste de Aceitação para IP dentro do Limite (TestRateLimiter_LimitByIP_Accept): Verifica se um IP que está dentro do limite de requisições é aceito.
    -> Teste de Aceitação para Token dentro do Limite (TestRateLimiter_LimitByToken_Accept): Verifica se um token que está dentro do limite de requisições é aceito.
    -> Teste de Sobreposição de Limite de IP por Token (TestRateLimiter_OverriddenByToken): Verifica se o limite de requisições do token sobrepõe o limite de requisições do IP.

## Como Modificar o Comportamento dos Testes

 Alterando o Comportamento do Mock: Se você deseja modificar o comportamento do rate limiter durante os testes, você pode alterar os mocks no arquivo limiter_test.go. O mock MockRateLimiter simula o comportamento do método Allow que é responsável por verificar se o limite foi excedido. Você pode modificar os valores retornados por ele para testar diferentes cenários de aceitação ou bloqueio.
 Exemplo de alteração:

```go 
    mockRateLimiter.On("Allow", mock.Anything, "ip:192.168.1.1", 5).Return(true, time.Duration(0)) // Permitir requisição
    mockRateLimiter.On("Allow", mock.Anything, "ip:192.168.1.1", 5).Return(false, time.Minute)     //  Bloquear requisição
```
## Alterando os Limites de Requisições: Você pode modificar o número de requisições permitidas por segundo diretamente nas chamadas do método Allow no código do teste. Por exemplo:

```go
    allowed, _ := limiter.Allow(context.Background(), "ip:192.168.1.1", 5) // Limite de 5 requisições por segundo
    allowed, _ = limiter.Allow(context.Background(), "token:ABC", 10) // Limite de 10 requisições por segundo
```
Modificando os Testes: Se você quiser adicionar novos cenários de teste, basta criar novas funções que seguem o mesmo padrão dos testes já implementados. Utilize o mock para simular os comportamentos esperados e as asserções para verificar os resultados.

## Exemplos de Testes

-> Teste de Limitação por IP
```go
    mockRateLimiter.On("Allow", mock.Anything, "ip:192.168.1.1", 5).Return(false, time.Minute)
    allowed, _ := limiter.Allow(context.Background(), "ip:192.168.1.1", 5)
    assert.False(t, allowed, "O IP deve ter excedido o limite de requisições")
```
-> Teste de Limitação por Token
```go
    mockRateLimiter.On("Allow", mock.Anything, "token:ABC", 10).Return(true, time.Duration(0))
    allowed, _ := limiter.Allow(context.Background(), "token:ABC", 10)
    assert.True(t, allowed, "O token deve estar dentro do limite de requisições")
```
-> Teste via Curl - Usando Token
```batch
curl -H "API_KEY: ABC" http://localhost:8080/
```
-> Teste via Curl - Usando IP
```batch
curl -H "X-Forwarded-For: 192.168.1.1" http://localhost:8080/
```
## Como Modificar o Comportamento do Rate Limiter
 -> Modificando o Comportamento do Rate Limiter:
O rate limiter utiliza uma estratégia baseada na interface RateLimiter. Para modificar o comportamento da limitação, você pode implementar diferentes estratégias de rate limiting, como:
    
   []- Token Bucket: Onde você tem um "balde" que se enche a cada intervalo de tempo e é consumido por cada requisição.
   []- Leaky Bucket: Onde as requisições são processadas em uma taxa constante, independentemente de quando são recebidas.

Basta criar uma nova implementação de RateLimiter e passá-la para a função NewRateLimiter durante a inicialização do servidor.

