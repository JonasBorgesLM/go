<!-- Jonas Borges L Moraes -->
<!-- jonasleo92@yahoo.com.br -->

# 🚀 Rate limiter

## Descrição

O _Rate Limiter_ é um projeto em Go que fornece um mecanismo para controlar a taxa de solicitações recebidas por um sistema ou API, agindo como middleware. Ele desempenha um papel crucial na prevenção de abusos, na proteção de recursos e na garantia de um uso equitativo. Este projeto utiliza o Redis para armazenar os dados do limite de taxa, mas é facilmente adaptável para utilizar outro sistema de armazenamento, desde que respeite a mesma interface.

## Funcionalidades

- **Funcionamento como Middleware:** Integra-se facilmente como middleware no servidor web.
- **Configuração da Taxa de Requisições:** Permite configurar de forma flexível a taxa máxima de requisições por segundo.
- **Opção de Bloqueio por Tempo:** Oferece funcionalidade para configurar o tempo de bloqueio do IP ou do Token caso o limite de requisições seja excedido.
- **Configuração via Variáveis de Ambiente ou .env:** As configurações de limite são facilmente configuráveis via variáveis de ambiente ou por meio de um arquivo ".env" na raiz do projeto.
- **Flexibilidade na Limitação por IP ou Token:** Pode ser configurado para limitar tanto por IP quanto por token de acesso.
- **Resposta Adequada ao Exceder o Limite:** Responde corretamente com o código de status HTTP 429 e uma mensagem informativa quando o limite de requisições é excedido.
- **Armazenamento e Consulta dos Dados no Redis:** Todas as informações do Limitador de Taxa são armazenadas e consultadas de forma eficiente em um banco de dados Redis, com a opção de utilizar docker-compose para facilitar a configuração do Redis.
- **Estratégia de Persistência Flexível:** Implementa uma estratégia que permite a fácil substituição do Redis por outro mecanismo de persistência, garantindo a flexibilidade do sistema.
- **Separação da Lógica do Middleware:** A lógica do limitador de taxa é modular e separada do middleware, facilitando a manutenção e a extensibilidade do código.

## Instalação

### Testes

Os testes de sobrecarga podem ser feitos através do Apache Bench.

`sudo apt install apache2-utils`

Use comandos para solicitar diversas requisições como o exemplo:

`ab -n 50 -c 1 -k -H  "API_KEY: JONAS"  "http://localhost:8080/"`

## Contribuições

Contribuições para o projeto Limitador de Taxa são bem-vindas! Sinta-se à vontade para enviar relatórios de bugs, solicitações de funcionalidades ou pull requests através do GitHub.

## Licença

Este projeto é licenciado sob a [Licença MIT](LICENSE).
