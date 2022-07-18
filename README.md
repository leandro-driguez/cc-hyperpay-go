# cc-hyperpay-go
Contrato inteligente para el manejo de cuentas con fondos.

`main.go` es el punto de entrada del contrato inteligente, el cual debe ser desplegado en la red [test-network-optativo-nanobash](https://github.com/ic-matcom/test-network-optativo-nanobash).

Implementamos una CLI para la comunicación con el contrato. Esta tiene como punto de entrada el archivo `client/main/main.go`. Para su correcto funcionamiento este repo debe ser clonado directamente dentro del repo [test-network-optativo-nanobash](https://github.com/ic-matcom/test-network-optativo-nanobash). Luego hacemos

```console
$ cd cc-hyperpay-go/client/main
$ go build main.go -o hyperpay
$ ./hyperpay
```

Al ejecutar el programa se muestra la ayuda de la aplicación, la cual expone los comandos disponibles. En la siguiente tabla se relacionan estos comandos con las funciones del contrato inteligente.

| Comando | Función en el cc | Ejemplo | Descripción |
|--------|--------|--------|--------|
| init | InitLedger | `./hyperpay init` | Coloca en la blockchain cuentas con IDs *account1*, *account2*, ..., *account5*. |
| read | ReadAccount | `./hyperpay read account1` | Consulta los datos de la cuenta con ID igual a *account1*. |
| exists | AccountExists | `./hyperpay exists account1` | Consulta la existencia en la blockchain de la cuenta con ID igual a *account1*. |
| delete | DeleteAccount | `./hyperpay delete account1` | Elimina la cuenta con ID igual a *account1*. |
| create | CreateAccount | `./hyperpay create new_account 120 BCC` | Crea una cuenta perteneciente al banco *BCC*, con ID igual a *new_account*, con saldo igual a 120. |
| transfer | Transfer | `./hyperpay transfer account1 account2 50` | Transfiere 50 del saldo de la cuenta con ID igual a *account1* a la cuenta con ID igual a *account2*. |
| txs | GetAllTxs | `./hyperpay txs account1` | Consulta todos los estados por los que ha transitado la cuenta con ID igual a *account1*. |