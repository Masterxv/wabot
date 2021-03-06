# Wabot

## Introdução

O objetivo desse projeto é enviar e receber mensagens utilizando o WhatsApp. Essa implementação é uma versão mais comercial do repositório go-whatsapp (https://github.com/Rhymen/go-whatsapp).

Desenvolvi essa solução pois a API oficial do WhatsApp é muito cara e restritiva.

## !!! Avisos !!!
Se seu objetivo é enviar cerca de 1.000 mensagens / min. , **não utilize essa api**. Ela foi feita para envios moderados, com baixa taxa de envio (cerca de 1 envio a cada 30-50 segundos) **por Chip**. 

Vale ressaltar que é possível a configuração de um ou mais Chips, podendo com isso **aumentar a taxa de envio**.


## Requisitos 

- GO
- MySQL 
- Um celular com conexão com internet
- Um chip com o WhatsApp Business configurado

## Configurações

- Execute o comando `go get` para instalar as dependências necessárias para a utilização da api.
- Crie um banco de dados `wabot`
- Execute o arquivo `/storage/wabot.sql` 
- Renomeie o arquivo `.env.example` para `.env` na pasta raiz do projeto
- Preencha as informações de conexão com o banco de dados
- Ainda no arquivo `.env`, preencha os endpoints:
    - Que fornecerá a fila de envio (`QUEUE_URL`)
    - Que removerá um envio da fila no servidor (`REMOVE_QUEUE_URL`)
    - Que receberá as respostas (`RESPONSES_URL`)
- Preencha as tabelas
    - `wabot_project` com o nome do seu projeto
    - `wabot_sender` com o número que está utilizando - o número do seu Chip
   
## Testando o envio

Após seguir os passos descritos anteriormente:

- Crie uma linha na tabela `wabot_queue` contendo 
    - `sender_id`
    - `número` que deseja enviar (pode ser o seu próprio número para fins de teste)
    - `send_date` e `send_time`: a data e hora que está agendado o envio. Preencha a data e hora atual para ocorrer o envio imediato
    - As demais colunas são opcionais
- Execute o arquivo main.go `go main.go`
    - Obs: Da primeira vez que rodar a aplicação irá exibir um QR Code na tela. Escaneie esse QR Code com seu WhatsApp
- Perceberá que é gerado um número aleatório como um `timeout` para o envio. Isso é para o detector de SPAM do WhatsApp não perceber atividade automatizada.
- Se seguiu todos os passos corretamente, o envio será feito com sucesso.
 
## Recomendações

- Utilize o WhatsApp Business.
- Não tente rodar esse App com fins de Marketing - Acredite, o WhatsApp irá te bloquear muito rápido.
- Para cada mensagem, tente variar o conteúdo. Dessa forma faz com que o WhatsApp não perceba atividade automatizada.
- Mantenha o celular que contém o CHIP **sempre** conectado no wi-fi.
- Preferencialmente, mantenha a tela do celular sempre ligada, conectado a um carregador e com o WhatsApp Business sempre aberto.
- É possível utilizar o WhatsApp Web no Chip que está enviando **apenas** quando os disparos não estão sendo feitos.


## Roadmap de novas funções / Contribuição

- Asegurar que a aplicação fique rodando 24/7.
- Criação de um Dashboard para controlar os envios / respostas
- Hospedar em um server
- Adaptar para rodar mais de uma instância com Chip de disparo diferente - para aumentar a taxa de envio
- Tentar encontrar um timeout menor que não seja bloqueado
- Desenvolver envio de mensagem baseado na resposta do usuário - ChatBot
- Adicionar função para BlackList de contatos

## Comparação com a api do WhatsApp Business

| Plataforma | Valor / envio | Taxa de habilitação | Risco de bloqueio do número | Permite Marketing | Modelo fixo de mensagem | Disparo / min. | Captura de resposta |  Envio deme imag | Permite envio de mais um número |
|:-----------------------------:|---------------|---------------------|----------------------------------|--------------------------------|-------------------------------|----------------|------------------------|------------------|-----------------------------------|
| WhatsApp Business API oficial | 0,36 | 6.000,00 | Não | Não | Sim | Ilimitado | Sim - por WebHook | Sim | Não - limitado a apenas um número |
| WABOT | 0,00 | 17,00 | Sim, se utilizar de forma errada | Sim - mas poderá ser bloqueado | Não - permite qualquer modelo | 2 | Sim - a cada 5 minutos | Não | Quantos números precisar |
