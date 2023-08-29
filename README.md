# bwg

Транзакциионная система, позволяющая класть и снимать деньги с баланса

## Методы: 
```put: 127.0.0.1:8080/output```

В body необходимо указать: userId, count, url

exemple:

![image](https://github.com/almazius/bwg/assets/101062396/6231bc27-2e14-43e2-b3c7-63f73ec5d141)


```put: 127.0.0.1:8080/input```

В body необходимо указать: userId, count

Если перейти по адресу, который пришел ответом транзакция подтвердится

exemple:

![image](https://github.com/almazius/bwg/assets/101062396/7ee16228-911d-44d7-a97c-cbfcd7d4035e)

