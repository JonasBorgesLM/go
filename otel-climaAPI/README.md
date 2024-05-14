<!-- Jonas Borges L Moraes -->
<!-- jonasleo92@yahoo.com.br -->

# 🚀 ClimaAPI

## Descrição

Este é um sistema desenvolvido em Go que recebe um CEP válido de 8 dígitos, identifica a cidade correspondente e retorna o clima atual em Celsius, Fahrenheit e Kelvin.

O sistema está sendo monitorado usando Open Telemetry e Zipkin.

## Requisitos

- O sistema deve receber um CEP válido de 8 dígitos.
- Deve realizar a pesquisa do CEP para encontrar o nome da localização e retornar as temperaturas formatadas em Celsius, Fahrenheit e Kelvin.
- Deve responder adequadamente nos seguintes cenários:
  - Em caso de sucesso:
    - Código HTTP: 200
    - Response Body: `{ "temp_C": 28.5, "temp_F": 28.5, "temp_K": 28.5 }`
  - Em caso de falha, caso o CEP não seja válido (com formato correto):
    - Código HTTP: 422
    - Mensagem: `invalid zipcode`
  - Em caso de falha, caso o CEP não seja encontrado:
    - Código HTTP: 404
    - Mensagem: `can not find zipcode`

## Endpoints

### GET /clima/{cep}

- **cep**: O CEP válido de 8 dígitos.

#### Exemplo de uso

```
GET /clima/12345678
```

#### Respostas

- **Sucesso**:

  - Código HTTP: 200
  - Body: `{ "temp_C": 28.5, "temp_F": 28.5, "temp_K": 28.5 }`

- **Falha - CEP inválido**:

  - Código HTTP: 422
  - Mensagem: `invalid zipcode`

- **Falha - CEP não encontrado**:
  - Código HTTP: 404
  - Mensagem: `can not find zipcode`

## Instalação local

- Clone o repositório.
- Execute `docker-compose up` para iniciar os serviços.

## Testes

Pode usar os arquivos .http do diretorio /http-test para testar a aplicação. Esses arquivos permitem que você defina e execute solicitações HTTP diretamente do seu editor de código.

### Host da aplicação local

`http://localhost:8080'

## Zipkin

`http://localhost:9411'

## Contribuições

Contribuições para o projeto Estresse são bem-vindas! Sinta-se à vontade para enviar relatórios de bugs, solicitações de funcionalidades ou pull requests através do GitHub.

## Licença

Este projeto é licenciado sob a [Licença MIT](LICENSE).
