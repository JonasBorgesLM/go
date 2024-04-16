<!-- Jonas Borges L Moraes -->
<!-- jonasleo92@yahoo.com.br -->

# 🚀 Teste de estresse - CLI

## Descrição

O _Teste de estresse - CLI_ é um sistema de linha de comando escrito em Go para realizar testes de carga em um serviço da web. Ele foi desenvolvido para permitir que os usuários testem a capacidade de resposta e a escalabilidade de seus serviços da web. Com este sistema, você pode especificar a URL do serviço, o número total de solicitações a serem enviadas e o nível de concorrência desejado para simular diferentes cenários de uso. É uma ferramenta útil para identificar gargalos de desempenho e otimizar a infraestrutura de servidores da web.

## Funcionalidades

- **Envio de Requests HTTP:** O sistema é capaz de enviar solicitações HTTP para uma URL específica, permitindo testar a capacidade de resposta do serviço da web.

- **Distribuição de Requests:** Os requests são distribuídos de acordo com o nível de concorrência definido pelo usuário, o que permite simular diferentes cargas de tráfego e identificar possíveis gargalos de desempenho.

- **Cumprimento do Número Total de Requests:** O sistema garante que o número total de requests especificado seja cumprido durante o teste, fornecendo resultados precisos sobre a capacidade de processamento do serviço da web.

- **Geração de Relatório:** Ao final dos testes, o sistema apresenta um relatório abrangente que inclui:
  - O tempo total gasto na execução dos testes.
  - A quantidade total de requests realizados.
  - A quantidade de requests com status HTTP 200.
  - A distribuição de outros códigos de status HTTP, como 404, 500, entre outros.

## Instalação

- Clone o repositório.

## Parameters

- `--url`: URL do serviço a ser testado.
- `--requests`: Número total de requests.
- `--concurrency`: Número de chamadas simultâneas.

### Testes

## Contribuições

Contribuições para o projeto Limitador de Taxa são bem-vindas! Sinta-se à vontade para enviar relatórios de bugs, solicitações de funcionalidades ou pull requests através do GitHub.

## Licença

Este projeto é licenciado sob a [Licença MIT](LICENSE).
