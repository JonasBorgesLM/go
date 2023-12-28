### Subir o banco e o rabbitmq:
`docker-compose up -d`

### Criar as tabelas do banco de dados:
`make migrate`

### Rodar a aplicação:
`go run ./main.go ./wire_gen.go`


### Testando os servidores:

#### Usar o graphql playground para testar o servidor GraphQL:
[Graphql - playground](http://localhost:8080)

#### Usar o evans para testar o servidor gRpc
`evans -r repl`

#### Usar os arquivos da pasta /api para testar o servidor web