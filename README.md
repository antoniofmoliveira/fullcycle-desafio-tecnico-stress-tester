# FullCycle - Desafio Técnico - Stress Test

## Solução:

### Projeto: DESAFIO_TECNICO_STRESS_TEST

* contém Makefile
  * make runtestserver - executa o server de testes localmente
  * make runstresstest - executa o stress test contra o server localmente

* subprojeto: server
  * `cmd/main.go` - servidor
    * devolve erros 500 e 429 randomicamente em 3% das chamadas
  * `api/request.http` - request http usando "REST Cliente"
  * `Dockerfile` - Dockerfile para construção da imagem para esse servidor

* subprojeto: stress-tester
  * `cmd/main.go` - trata flags de entrada e executa o stress test
  * `internal` - pacotes internos do app
    * `db` - banco de dados em memória sqllite3
    * `dto` - modelos de dados transferidos entre camadas
    * `entity` - entidades do domínio
    * `pool` - pool de httoclient e banco de dados
      * `db-pool` - pool de banco de dados
      * `htt-client-pool` - pool de httpclient para envio de grande volume de requests
    * `report` - gerador de relatório
    * `stats` - calculos estatísticos
    * `usecase` - usecase para execução do stress test
      * calcular quantos rounds serao feitos e quantos requests concorrentes serao feitos em cada round
      * executa os requests
      * imprime o relatório
  * `Dockerfile` - Dockerfile para construção da imagem para esse stress-tester
    * observar que depende de bibliotecas C por causa do driver sqlite
    * CMD contém parâmetros padrão para o stress-tester

* Relatório:
  * relação de rounds executados. observar que o último round pode ter menos request em virtude da divisão inteira deixar restos
  * total de requests e tempo de execução
  * resumo
    * quantidade de requests
    * quantidade de requests com erro
    * tempo médio de resposta
    * menor tempo de resposta
    * maior tempo de resposta
    * quantidade de erros de rede (server não respondeu ao request)
  * resumo por status code
    * erro -1 indica erros de rede (server não respondeu ao request)
  * distribuição de erros por percentil de tempo das respostas (10%, 25%, 50%, 75%, 90%, 99%) - quão distante estão os piores tempos dos melhores tempos

```bash
Round  0 Running  10  requests for endpoint  http://localhost:8080
Round  1 Running  10  requests for endpoint  http://localhost:8080
Round  2 Running  10  requests for endpoint  http://localhost:8080
Round  3 Running  10  requests for endpoint  http://localhost:8080
Round  4 Running  10  requests for endpoint  http://localhost:8080
Round  5 Running  10  requests for endpoint  http://localhost:8080
Round  6 Running  10  requests for endpoint  http://localhost:8080
Round  7 Running  10  requests for endpoint  http://localhost:8080
Round  8 Running  10  requests for endpoint  http://localhost:8080
Round  9 Running  10  requests for endpoint  http://localhost:8080
Round  10 Running  5  requests for endpoint  http://localhost:8080
Finished  105  requests for endpoint  http://localhost:8080  in  15.88441ms
      Rate           Error        Avg Time        Min Time        Max Time       Net Error
       105               7      11.990396ms     11.079686ms     13.384475ms              0

Status  # Responses
200             98
429              4
500              3

Percentile        Duration
P10             11.563215ms
P25             11.919995ms
P50             12.076606ms
P75             12.309575ms
P90             12.718199ms
P99             13.30991ms
```

### Execução

* no raiz do projeto execute `make build` para gerar docker image do stress-tester
* execute `docker run stresstester  --url=http://google.com --requests=105 --concurrency=10` para ver o relatório gerado

#### Execução no Docker

* no raiz do projeto execute `make run-server` para executar o server de exemplo no docker
* execute `docker inspect server | grep "IPAddress"` para obter o IP do server no docker
* execute `docker run stresstester  --url=http://172.17.0.2:8080 --requests=105 --concurrency=10` para ver o relatório gerado. substitua `172.17.0.2` pelo IP do server no docker obtido no passo anterior.

## Objetivo:

Criar um sistema CLI em Go para realizar testes de carga em um serviço web. O usuário deverá fornecer a URL do serviço, o número total de requests e a quantidade de chamadas simultâneas.

O sistema deverá gerar um relatório com informações específicas após a execução dos testes.

## Entrada de Parâmetros via CLI:

--url: URL do serviço a ser testado.<br>
--requests: Número total de requests.<br>
--concurrency: Número de chamadas simultâneas.

# Execução do Teste:

* Realizar requests HTTP para a URL especificada.
* Distribuir os requests de acordo com o nível de concorrência definido.
* Garantir que o número total de requests seja cumprido.

# Geração de Relatório:

* Apresentar um relatório ao final dos testes contendo:
  * Tempo total gasto na execução
  * Quantidade total de requests realizados.
  * Quantidade de requests com status HTTP 200.
  * Distribuição de outros códigos de status HTTP (como 404, 500, etc.).

1. Execução da aplicação:

* Poderemos utilizar essa aplicação fazendo uma chamada via docker. Ex:
  * `docker run <sua imagem docker> —url=http://google.com —requests=1000 —concurrency=10`
