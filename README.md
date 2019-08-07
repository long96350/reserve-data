# Data fetcher for KyberNetwork reserve
[![Go Report Card](https://goreportcard.com/badge/github.com/KyberNetwork/reserve-data)](https://goreportcard.com/report/github.com/KyberNetwork/reserve-data)
[![Build Status](https://travis-ci.org/KyberNetwork/reserve-data.svg?branch=develop)](https://travis-ci.org/KyberNetwork/reserve-data)

This repo is contains two components:

- core: 
	- interacts with blockchain to get/set rates for tokens pair
	- buy/sell with centralized exchanges (binance, huobi, etc)
(For more detail, take a look to interface ReserveCore in intefaces.go)

- stat:  
	- fetch tradelogs from blockchain and do aggregation and save its data to database and allow client to query

(For more detail, find ReserveStat interface in interfaces.go)

## Compile it

```shell
cd cmd && go build -v
```

a `cmd` executable file will be created in `cmd` module.

## Run the reserve data

1. You need to prepare a `config.json` file inside `cmd` module. The file is described in later section.
2. You need to prepare a JSON keystore file inside `cmd` module. It is the keystore for the reserve owner.
3. Make sure your working directory is `cmd`. Run `KYBER_EXCHANGES=binance,huobi ./cmd` in dev mode.

### Manual

```shell
cd cmd
```

- Run core only

```shell
KYBER_EXCHANGES="binance,huobi" KYBER_ENV=production ./cmd server --log-to-stdout
```

- Run stat only

```shell
KYBER_ENV=production ./cmd server --log-to-stdout --enable-stat --no-core
```

### Docker (recommended)

This repository will build docker images and public on [docker hub](https://hub.docker.com/r/kybernetwork/reserve-data/tags/), you can pull image from docker hub and run:

- Run core only

```shell
docker run -p 8000:8000 -v /location/of/config.json:/go/src/github.com/KyberNetwork/reserve-data/cmd/config.json -e KYBER_EXCHANGES="binance,huobi" KYBER_ENV="production" kybernetwork/reserve-data:develop server --log-to-stdout
```

- Run stat only 

```shell
docker run -p 8000:8000 -v /location/of/config.json:/go/src/github.com/KyberNetwork/reserve-data/cmd/config.json -e KYBER_ENV="production" kybernetwork/reserve-data:develop server --enable-stat --no-core --log-to-stdout
```

**Note** : 

- KYBER_ENV includes "dev, simulation and production", different environment mode uses different settings (check cmd folder for settings file).  

- reserve-data requires config.json file to run, so you need to -v (mount config.json file to docker) so it can run.

## Config file

sample:

```json
{
  "binance_key": "your binance key",
  "binance_secret": "your binance secret",
  "huobi_key": "your huobi key",
  "huobi_secret_key": "your huobi secret",
  "kn_secret": "secret key for people to sign their requests to our apis. It is ignored in dev mode.",
  "kn_readonly": "read only key for people to sign their requests, this key can read everything but cannot execute anything",
  "kn_configuration": "key for people to sign their requests, this key can read everything and set configuration such as target quantity",
  "kn_confirm_configuration": "key for people to sign ther requests, this key can read everything and confirm target quantity, enable/disable setrate or rebalance",
  "keystore_path": "path to the JSON keystore file, recommended to be absolute path",
  "passphrase": "passphrase to unlock the JSON keystore",
  "keystore_deposit_path": "path to the JSON keystore file that will be used to deposit",
  "passphrase_deposit": "passphrase to unlock the JSON keystore",
  "keystore_intermediator_path": "path to JSON keystore file that will be used to deposit to Huobi",
  "passphrase_intermediate_account": "passphrase to unlock JSON keystore",
  "aws_access_key_id": "your aws key ID",
  "aws_secret_access_key": "your aws scret key",
  "aws_expired_stat_data_bucket_name" : "AWS bucket for expired stat data (already created)",
  "aws_expired_reserve_data_bucket_name" : "AWS bucket for expired reserve data (already created)",
  "aws_log_bucket_name" :"AWS bucket for log backup(already created)",
  "aws_region":"AWS region"
}
```

## APIs

### Get token rates from blockchain

```shell
<host>:8000/getrates
```

eg:

```shell
curl -X GET "http://127.0.0.1:8000/getrates"
```
response:

```json
{
    "data": {
        "MCO": {
            "Valid": true,
            "Error": "",
            "Timestamp": "1515412582435",
            "ReturnTime": "1515412582710",
            "BaseBuy": 63.99319226272073,
            "CompactBuy": 21,
            "BaseSell": 0.014716371218820246,
            "CompactSell": -20,
            "Rate": 0,
            "Block": 2420849
        },
        "OMG": {
            "Valid": true,
            "Error": "",
            "Timestamp": "1515412582435",
            "ReturnTime": "1515412582710",
            "BaseBuy": 44.45707162223901,
            "CompactBuy": 30,
            "BaseSell": 0.021183301968644246,
            "CompactSell": -29,
            "Rate": 0,
            "Block": 2420849
        },
        "PAY": {
            "Valid": true,
            "Error": "",
            "Timestamp": "1515412582435",
            "ReturnTime": "1515412582710",
            "BaseBuy": 295.08854913901575,
            "CompactBuy": -13,
            "BaseSell": 0.003191406699999999,
            "CompactSell": 13,
            "Rate": 0,
            "Block": 2420849
        }
    },
    "success": true,
    "timestamp": "1515412583215",
    "version": 1515412582435
}
```


### Get all token rates from blockchain

```shell
<host>:8000/get-all-rates
```

url params:
*fromTime*: optional, get all rates from this timepoint (millisecond)
*toTime*: optional, get all rates to this timepoint (millisecond)

eg:

```shell
curl -X GET "http://127.0.0.1:8000/get-all-rates"
```

response

```json
{
    "data": [
        {
            "Version": 0,
            "Valid": true,
            "Error": "",
            "Timestamp": "1517280618739",
            "ReturnTime": "1517280619071",
            "Data": {
                "SNT": {
                    "Valid": true,
                    "Error": "",
                    "Timestamp": "1517280618739",
                    "ReturnTime": "1517280619071",
                    "BaseBuy": 4053.2170631085987,
                    "CompactBuy": 43,
                    "BaseSell": 0.000233599514875301,
                    "CompactSell": -3,
                    "Rate": 0,
                    "Block": 5635245
                }
            }
        },
        {
            "Version": 0,
            "Valid": true,
            "Error": "",
            "Timestamp": "1517280621738",
            "ReturnTime": "1517280622251",
            "Data": {
                "EOS": {
                    "Valid": true,
                    "Error": "",
                    "Timestamp": "1517280621738",
                    "ReturnTime": "1517280622251",
                    "BaseBuy": 87.21360760013062,
                    "CompactBuy": 0,
                    "BaseSell": 0.0128686459657361,
                    "CompactSell": 0,
                    "Rate": 0,
                    "Block": 5635245
                },
                "ETH": {
                    "Valid": true,
                    "Error": "",
                    "Timestamp": "1517280621738",
                    "ReturnTime": "1517280622251",
                    "BaseBuy": 0,
                    "CompactBuy": 32,
                    "BaseSell": 0,
                    "CompactSell": -14,
                    "Rate": 0,
                    "Block": 5635245
                },
                "KNC": {
                    "Valid": true,
                    "Error": "",
                    "Timestamp": "1517280621738",
                    "ReturnTime": "1517280622251",
                    "BaseBuy": 307.05930436561505,
                    "CompactBuy": -34,
                    "BaseSell": 0.003084981280661941,
                    "CompactSell": 81,
                    "Rate": 0,
                    "Block": 5635245
                },
                "OMG": {
                    "Valid": true,
                    "Error": "",
                    "Timestamp": "1517280621738",
                    "ReturnTime": "1517280622251",
                    "BaseBuy": 65.0580993582104,
                    "CompactBuy": 32,
                    "BaseSell": 0.014925950060437398,
                    "CompactSell": -14,
                    "Rate": 0,
                    "Block": 5635245
                },
                "SALT": {
                    "Valid": true,
                    "Error": "",
                    "Timestamp": "1517280621738",
                    "ReturnTime": "1517280622251",
                    "BaseBuy": 152.3016783627643,
                    "CompactBuy": 9,
                    "BaseSell": 0.006196212698403499,
                    "CompactSell": 23,
                    "Rate": 0,
                    "Block": 5635245
                },
                "SNT": {
                    "Valid": true,
                    "Error": "",
                    "Timestamp": "1517280621738",
                    "ReturnTime": "1517280622251",
                    "BaseBuy": 4053.2170631085987,
                    "CompactBuy": 43,
                    "BaseSell": 0.000233599514875301,
                    "CompactSell": -3,
                    "Rate": 0,
                    "Block": 5635245
                }
            }
        }
    ],
    "success": true
}
```


### Get trade history for an account (signing required)

```shell
  <host>:8000/tradehistory  
  params:
  - fromTime: millisecond (required)
  - toTime: millisecond (required)
  Restriction: toTime - fromTime <= 3 days (in millisecond)
```

eg:

```shell
curl -X GET "http://localhost:8000/tradehistoryfromTime=1516116380102&toTime=18446737278344972745"
```

response:

```json
{"data":{"Version":1517298257114,"Valid":true,"Timestamp":"1517298257115","Data":{"binance":{"EOS-ETH":[],"KNC-ETH":[{"ID":"548002","Price":0.003038,"Qty":50,"Type":"buy","Timestamp":1516116380102},{"ID":"548003","Price":0.0030384,"Qty":7,"Type":"buy","Timestamp":1516116380102},{"ID":"548004","Price":0.003043,"Qty":16,"Type":"buy","Timestamp":1516116380102},{"ID":"548005","Price":0.0030604,"Qty":29,"Type":"buy","Timestamp":1516116380102},{"ID":"548006","Price":0.003065,"Qty":29,"Type":"buy","Timestamp":1516116380102},{"ID":"548007","Price":0.003065,"Qty":130,"Type":"buy","Timestamp":1516116380102}],"OMG-ETH":[{"ID":"123980","Price":0.020473,"Qty":48,"Type":"buy","Timestamp":1512395498231},{"ID":"130518","Price":0.021022,"Qty":13.49,"Type":"buy","Timestamp":1512564108827},{"ID":"130706","Price":0.020202,"Qty":9.93,"Type":"sell","Timestamp":1512569059460},{"ID":"140078","Price":0.019098,"Qty":11.07,"Type":"buy","Timestamp":1512714826339},{"ID":"140157","Price":0.019053,"Qty":7.68,"Type":"sell","Timestamp":1512716338997},{"ID":"295923","Price":0.020446,"Qty":4,"Type":"buy","Timestamp":1514360742162}],"SALT-ETH":[],"SNT-ETH":[]},"bittrex":{"OMG-ETH":[{"ID":"eb948865-6261-4991-8615-b36c8ccd1256","Price":0.01822057,"Qty":1,"Type":"buy","Timestamp":18446737278344972745}],"SALT-ETH":[],"SNT-ETH":[]}}},"success":true}
```



### Get exchange balances, reserve balances, pending activities at once (signing required)

```shell
<host>:8000/authdata
```

eg:

```shell
curl -X GET "http://localhost:8000/authdata"
```

response:

```json
{"data":{"Valid":true,"Error":"","Timestamp":"1514114408227","ReturnTime":"1514114408810","ExchangeBalances":{"bittrex":{"Valid":true,"Error":"","Timestamp":"1514114408226","ReturnTime":"1514114408461","AvailableBalance":{"ETH":0.10704306,"OMG":2.97381136},"LockedBalance":{"ETH":0,"OMG":0},"DepositBalance":{"ETH":0,"OMG":0}}},"ReserveBalances":{"ADX":{"Valid":true,"Error":"","Timestamp":"1514114408461","ReturnTime":"1514114408799","Balance":0},"BAT":{"Valid":true,"Error":"","Timestamp":"1514114408461","ReturnTime":"1514114408799","Balance":0},"CVC":{"Valid":true,"Error":"","Timestamp":"1514114408461","ReturnTime":"1514114408799","Balance":0},"DGD":{"Valid":true,"Error":"","Timestamp":"1514114408461","ReturnTime":"1514114408799","Balance":0},"EOS":{"Valid":true,"Error":"","Timestamp":"1514114408461","ReturnTime":"1514114408799","Balance":0},"ETH":{"Valid":true,"Error":"","Timestamp":"1514114408461","ReturnTime":"1514114408799","Balance":360169992138038352},"FUN":{"Valid":true,"Error":"","Timestamp":"1514114408461","ReturnTime":"1514114408799","Balance":0},"GNT":{"Valid":true,"Error":"","Timestamp":"1514114408461","ReturnTime":"1514114408799","Balance":0},"KNC":{"Valid":true,"Error":"","Timestamp":"1514114408461","ReturnTime":"1514114408799","Balance":0},"LINK":{"Valid":true,"Error":"","Timestamp":"1514114408461","ReturnTime":"1514114408799","Balance":0},"MCO":{"Valid":true,"Error":"","Timestamp":"1514114408461","ReturnTime":"1514114408799","Balance":0},"OMG":{"Valid":true,"Error":"","Timestamp":"1514114408461","ReturnTime":"1514114408799","Balance":23818094310417195708},"PAY":{"Valid":true,"Error":"","Timestamp":"1514114408461","ReturnTime":"1514114408799","Balance":0}},"PendingActivities":[]},"block": 2345678, "success":true,"timestamp":"1514114409088","version":39}
```

### Get all activityes (signing required)
```
<host>:8000/activities
GET request
url params: 
  fromTime: from timepoint - uint64, unix millisecond (optional if empty then get from first activity)
  toTime: to timepoint - uint64, unix millisecond (optional if empty then get to last activity)
```
Note: `fromTime` and `toTime` shouldn't be included into signing message.
### Get immediate pending activities (signing required)
```
<host>:8000/immediate-pending-activities
GET request
```

### Store processed data (signing required)
```
<host>:8000/metrics
POST request
form params:
  - timestamp: uint64, unix millisecond
  - data: string, in format of <token>_afpmid_spread|<token>_afpmid_spread|..., eg. OMG_0.4_5|KNC_1_2
```

### Get processed data (signing required)
```
<host>:8000/metrics
GET request
url params:
  - tokens: string, list of tokens to get data about, in format of <token_id>-<token_id>..., eg. OMG_DGD_KNC
  - from: uint64, unix millisecond
  - to: uint64, unix millisecond
```

response:
```json
{
    "data": {
        "DGD": [
            {
                "Timestamp": 19,
                "AfpMid": 4,
                "Spread": 5
            }
        ],
        "OMG": [
            {
                "Timestamp": 19,
                "AfpMid": 0.9,
                "Spread": 1
            }
        ]
    },
    "returnTime": 1514966512560,
    "success": true,
    "timestamp": 1514966512549
}
```
Returned data will only include datas that have timestamp in range of `[from, to]`


### Get pending pwis equation (signing required)
```
<host>:8000/pending-pwis-equation
GET request
```

response:
```
  {
    "success": true,
    "data":{"id":1517396850670,"data":"EOS_750_500_0.25|ETH_750_500_0.25|KNC_750_500_0.25|OMG_750_500_0.25|SALT_750_500_0.25"}
  }
```

### Get pwis equation (signing required)
```
<host>:8000/pwis-equation
GET request
```

response:
```
  {
    "success": true,
    "data":{"id":1517396850670,"data":"EOS_750_500_0.25|ETH_750_500_0.25|KNC_750_500_0.25|OMG_750_500_0.25|SALT_750_500_0.25"}
  }
```

### Set pwis equation (signing required)
```
<host>:8000/set-pwis-equation
POST request
form params:
  - data: required, string, must sort by token id by ascending order
  - id: optional, required to confirm target quantity
```
eg:
```
curl -X POST \
  http://localhost:8000/set-pwis-equation \
  -H 'content-type: multipart/form-data' \
  -F data= EOS_750_500_0.25|ETH_750_500_0.25|KNC_750_500_0.25|OMG_750_500_0.25|SALT_750_500_0.25 \
  -F id=1517396850670
```
response
```
  {
    "success": true,
  }
```

### Confirm pwis equation (signing required)
```
<host>:8000/confirm-pwis-equation
POST request
form params:
  - data: required, string, must sort by token id by ascending order
```
eg:
```
curl -X POST \
  http://localhost:8000/confirm-pwis-equation \
  -H 'content-type: multipart/form-data' \
  -F data=EOS_750_500_0.25|ETH_750_500_0.25|KNC_750_500_0.25|OMG_750_500_0.25|SALT_750_500_0.25
```
response
```
  {
    "success": true,
  }
```

### Reject pwis equation (signing required)
```
<host>:8000/reject-pwis-equation
POST request
```
eg:
```
curl -X POST \
  http://localhost:8000/reject-pwis-equation \
  -H 'content-type: multipart/form-data' \
```
response
```
  {
    "success": true,
  }
```

### Get trade logs

```
<host>:8000/tradelogs
GET request
```
response
```
  {
    "success": true
  }
```


### Get Cap by addr
```
<host>:8000/cap-by-address/:addr
GET request

Url params:
  - address (required): ethereum Address of the user. 

Return : a record of type UserCap struct {
	UserID     string
	Category   string
	DailyLimit float64
	TxLimit    float64
	Type       string}
  - if the address is Kyced, the kyced Cap is returned
  - otherwise return normal cap
```

### Get Cap by addr
```
<host>:8000/richguy/:addr
GET request

Url params:
  - address (required): ethereum Address of the user. 

Return : boolean
  - true if the address trading volume has pasted the limit, false otherwise
```
### Get User Lists
```
<host>:8000/get-user-list
GET request

Url params:
  - fromTime (millisecond - required): from time stamp
  - toTime (millisecond - required: to time stamp
Return : a list of records available with timestamp in between fromTime toTime, of type UserInfo struct {
	Addr      string  `json:"user_address"`
	Email     string  `json:"email"`
	ETHVolume float64 `json:"total_eth_volume"`
	USDVolume float64 `json:"total_usd_volume"`
}
```

### Get asset volume for aggregate time (hour, day, month)
```
<host>:8000/get-asset-volume
GET request

Url params:
  - fromTime (millisecond - required): from time stamp
  - toTime (millisecond - required: to time stamp
  - freq (required): frequency to get data (H/D/M)
  - asset (required): asset name (eg: ETH)
```

```
curl -x GET http://localhost:8000/get-asset-volume?fromTime=1520640035000&toTime=1520722835000&asset=eth&freq=M
```
response
```
  {"data":{"1520652360000":{"usd_amount":0.734518,"volume":0.001},"1520654280000":{"usd_amount":0.7297319999999999,"volume":0.001},"1520654820000":{"usd_amount":1.4581552500230603,"volume":0.001998206533389053},"1520656440000":{"usd_amount":0.7297319999999999,"volume":0.001},"1520656500000":{"usd_amount":0.7297319999999999,"volume":0.001},"1520656560000":{"usd_amount":0.7297319999999999,"volume":0.001}},"success":true}
```

### Get burn fee follow aggregate time (hour, day, month)
```
<host>:8000/get-burn-fee
GET request

Url params:
  - fromTime (millisecond - required): from time stamp
  - toTime (millisecond - required: to time stamp
  - freq (required): frequency to get data (H/D/M)
  - reserveAddr (required): reserve address to get burn fee
```

```
curl -x GET http://localhost:8000/get-burn-fee?fromTime=1520640035000&toTime=1520722835000&reserveAddr=0x2c5a182d280eeb5824377b98cd74871f78d6b8bc&freq=H
```
response
```
  {"data":{"1520650800000":0.00225,"1520654400000":0.005622982350062684},"success":true}
```

### Get wallet fee follow aggregate time (hour, day, month)
```
<host>:8000/get-wallet-fee
GET request

Url params:
  - fromTime (millisecond - required): from time stamp
  - toTime (millisecond - required: to time stamp
  - freq (required): frequency to get data (H/D/M)
  - reserveAddr (required): reserve address
  - walletAddr (required): wallet address to get fee
```

```
curl -x GET http://localhost:8000/get-wallet-fee?fromTime=1520640035000&toTime=1520722835000&reserveAddr=0x2c5a182d280eeb5824377b98cd74871f78d6b8bc&walletAddr=0x0000000000000000000000000000000000000000&freq=H
```
response
```
  {"data":{"1520650800000":0,"1520654400000":0},"success":true}
```


### Get user volume follow aggregate time (hour, day, month)
```
<host>:8000/get-user-volume
GET request

Url params:
  - fromTime (millisecond - required): from time stamp
  - toTime (millisecond - required: to time stamp
  - freq (required): frequency to get data (H/D/M)
  - userAddr (required): user address to get volume
```

```
curl -x GET http://localhost:8000/get-user-volume?fromTime=1520640035000&toTime=1520722835000&userAddr=0x8fa07f46353a2b17e92645592a94a0fc1ceb783f&freq=H
```
response
```
  {"data":{"1520650800000":0.734518,"1520654400000":0},"success":true}  
```

### Get rate from blockchain follow reserve (including sanity rate)
```
<host>:8000/get-reserve-rate
GET request

URL params:
  - fromTime (millisecond - required): from timestamp
  - toTime (millisecond - required): to timestamp
  - reserveAddr (required): Address of the reserve to get rate from
```

```
curl -x GET http://localhost:8000/get-reserve-rate?fromTime=1520650426000&reserveAddr=0x2C5a182d280EeB5824377B98CD74871f78d6b8BC
```

response
```
{"data":[{"Timestamp":0,"ReturnTime":1520655211398,"BlockNumber":5228238,"Data":{"APPC-ETH":{"ReserveRate":0.008393501685222925,"SanityRate":0.009476954807692308},"BAT-ETH":{"ReserveRate":0.004239837479770336,"SanityRate":0.004723026},"BQX-ETH":{"ReserveRate":0.000584106942517358,"SanityRate":0.000652623583333333},"ELF-ETH":{"ReserveRate":0.000111035861616385,"SanityRate":0.000123576933333333},"ENG-ETH":{"ReserveRate":0.000596961062855617,"SanityRate":0.000671348333333333},"EOS-ETH":{"ReserveRate":0.002752586323625439,"SanityRate":0.0029518775},"ETH-APPC":{"ReserveRate":117.09352189952426,"SanityRate":127.67814393478591},"ETH-BAT":{"ReserveRate":229.70817443088293,"SanityRate":256.19168727845243},"ETH-BQX":{"ReserveRate":1674.1030165099337,"SanityRate":1854.0549727299406},"ETH-ELF":{"ReserveRate":8741.854159397268,"SanityRate":9791.471331758756},"ETH-ENG":{"ReserveRate":1615.634600168794,"SanityRate":1802.3430459597466},"ETH-EOS":{"ReserveRate":356.63560217559376,"SanityRate":409.9086090123997},"ETH-GTO":{"ReserveRate":377.41338205276884,"SanityRate":432.7020247561754},"ETH-KNC":{"ReserveRate":3343.445798727388,"SanityRate":3740.5264791019335},"ETH-MANA":{"ReserveRate":2653.110602592891,"SanityRate":2961.048749629869},"ETH-OMG":{"ReserveRate":221.03211631662654,"SanityRate":247.71892159696137},"ETH-POWR":{"ReserveRate":2625.724091042635,"SanityRate":2849.3656923419003},"ETH-RDN":{"ReserveRate":49.46371879742714,"SanityRate":54.91076347236177},"ETH-REQ":{"ReserveRate":5123.294220111987,"SanityRate":5665.576472406068},"ETH-SALT":{"ReserveRate":532.7984920557698,"SanityRate":611.1450636146453},"ETH-SNT":{"ReserveRate":924.3533982883454,"SanityRate":1052.547224086108},"GTO-ETH":{"ReserveRate":0.002590177384056749,"SanityRate":0.002796381645502645},"KNC-ETH":{"ReserveRate":0.0002894500017402,"SanityRate":0.000323483875},"MANA-ETH":{"ReserveRate":0.00036813437957934,"SanityRate":0.000408639},"OMG-ETH":{"ReserveRate":0.004383746019560721,"SanityRate":0.004884568333333332},"POWR-ETH":{"ReserveRate":0.000369936605210205,"SanityRate":0.000424655916666666},"RDN-ETH":{"ReserveRate":0.01987031936393942,"SanityRate":0.02203575261904761},"REQ-ETH":{"ReserveRate":0.000191920526182855,"SanityRate":0.0002135705},"SALT-ETH":{"ReserveRate":0.001821407074188612,"SanityRate":0.00197989},"SNT-ETH":{"ReserveRate":0.001042608668464954,"SanityRate":0.001149592125}}},{"Timestamp":0,"ReturnTime":1520655227886,"BlockNumber":5228239,"Data":{"APPC-ETH":{"ReserveRate":0.000369936605210205,"SanityRate":0.000424655916666666},"BAT-ETH":{"ReserveRate":0.0002894500017402,"SanityRate":0.000323483875},"BQX-ETH":{"ReserveRate":0.002590177384056749,"SanityRate":0.002796381645502645},"ELF-ETH":{"ReserveRate":0.01987031936393942,"SanityRate":0.02203575261904761},"ENG-ETH":{"ReserveRate":0.000584106942517358,"SanityRate":0.000652623583333333},"EOS-ETH":{"ReserveRate":0.000191920526182855,"SanityRate":0.0002135705},"ETH-APPC":{"ReserveRate":2625.724091042635,"SanityRate":2849.3656923419003},"ETH-BAT":{"ReserveRate":3343.445798727388,"SanityRate":3740.5264791019335},"ETH-BQX":{"ReserveRate":377.41338205276884,"SanityRate":432.7020247561754},"ETH-ELF":{"ReserveRate":49.46371879742714,"SanityRate":54.91076347236177},"ETH-ENG":{"ReserveRate":1674.1030165099337,"SanityRate":1854.0549727299406},"ETH-EOS":{"ReserveRate":5123.294220111987,"SanityRate":5665.576472406068},"ETH-GTO":{"ReserveRate":229.70817443088293,"SanityRate":256.19168727845243},"ETH-KNC":{"ReserveRate":1615.634600168794,"SanityRate":1802.3430459597466},"ETH-MANA":{"ReserveRate":221.03211631662654,"SanityRate":247.71892159696137},"ETH-OMG":{"ReserveRate":8741.854159397268,"SanityRate":9791.471331758756},"ETH-POWR":{"ReserveRate":924.3533982883454,"SanityRate":1052.547224086108},"ETH-RDN":{"ReserveRate":532.7984920557698,"SanityRate":611.1450636146453},"ETH-REQ":{"ReserveRate":117.09352189952426,"SanityRate":127.67814393478591},"ETH-SALT":{"ReserveRate":2653.110602592891,"SanityRate":2961.048749629869},"ETH-SNT":{"ReserveRate":356.63560217559376,"SanityRate":409.9086090123997},"GTO-ETH":{"ReserveRate":0.004239837479770336,"SanityRate":0.004723026},"KNC-ETH":{"ReserveRate":0.000596961062855617,"SanityRate":0.000671348333333333},"MANA-ETH":{"ReserveRate":0.004383746019560721,"SanityRate":0.004884568333333332},"OMG-ETH":{"ReserveRate":0.000111035861616385,"SanityRate":0.000123576933333333},"POWR-ETH":{"ReserveRate":0.001042608668464954,"SanityRate":0.001149592125},"RDN-ETH":{"ReserveRate":0.001821407074188612,"SanityRate":0.00197989},"REQ-ETH":{"ReserveRate":0.008393501685222925,"SanityRate":0.009476954807692308},"SALT-ETH":{"ReserveRate":0.00036813437957934,"SanityRate":0.000408639},"SNT-ETH":{"ReserveRate":0.002752586323625439,"SanityRate":0.0029518775}}}],"success":true}
```


### Get trade summary follow timeframe (day)

```shell
<host>:8000/get-trade-summary
GET request

Url params:
  - fromTime (millisecond - required): from time stamp
  - toTime (millisecond - required): to time stamp  
  - timeZone (in range [-12,14], default to 0): the integer specific which UTC timezone to query
```

eg:

```shell
curl -x GET http://localhost:8000/get-trade-summary?fromTime=1519297149000&toTime=1519815549000
```

response

```json
{"data":{"1519344000000":{"eth_per_trade":0.55402703087424,"kyced_addresses":0,"new_unique_addresses":35,"total_burn_fee":0,"total_eth_volume":44.3221624699392,"total_trade":80,"total_usd_amount":30981.281202536768,"unique_addresses":50,"usd_per_trade":387.26601503170957},"1519430400000":{"eth_per_trade":0.17008867987348247,"kyced_addresses":0,"new_unique_addresses":17,"total_burn_fee":0,"total_eth_volume":8.674522673547607,"total_trade":51,"total_usd_amount":6060.828270348999,"unique_addresses":29,"usd_per_trade":118.83977000684311},"1519516800000":{"eth_per_trade":0.14234886960871,"kyced_addresses":0,"new_unique_addresses":9,"total_burn_fee":1.1025,"total_eth_volume":5.40925704513098,"total_trade":38,"total_usd_amount":3779.4100326337,"unique_addresses":18,"usd_per_trade":99.45815875351843},"1519603200000":{"eth_per_trade":0.5430574166436676,"kyced_addresses":0,"new_unique_addresses":39,"total_burn_fee":42.85336706164196,"total_eth_volume":45.07376558142441,"total_trade":83,"total_usd_amount":31497.3427579499,"unique_addresses":56,"usd_per_trade":379.4860573246976},"1519689600000":{"eth_per_trade":0.6014134385918366,"kyced_addresses":0,"new_unique_addresses":69,"total_burn_fee":79.03472646631772,"total_eth_volume":78.7851604555306,"total_trade":131,"total_usd_amount":55076.026979006005,"unique_addresses":92,"usd_per_trade":420.4276868626413},"1519776000000":{"eth_per_trade":0.40083191776618454,"kyced_addresses":0,"new_unique_addresses":64,"total_burn_fee":48.899026261678536,"total_eth_volume":52.50898122737018,"total_trade":131,"total_usd_amount":36662.138255818456,"unique_addresses":94,"usd_per_trade":279.8636508077745}},"success":true}
```

### Get a specific wallet's stats summary follow timeframe (day)

```shell
<host>:8000/get-wallet-stats
GET request

Url params:
  - fromTime (millisecond - required): from time stamp
  - toTime (millisecond - required): to time stamp  
  - timeZone (in range [-12,14], default to 0): the integer specific which UTC timezone to query
  - walletAddr (hex string - required) : to specific which wallet Address to query data from. It must be larger than 2^128 to be valid.
```

```shell
curl -x GET http://localhost:8000/get-wallet-stats?fromTime=1521914061000&toTime=1523000461000&walletAddr=0xb9e29984fe50602e7a619662ebed4f90d93824c7
```
response

```json
{"data":{"1521936000000":{"eth_per_trade":0.15169175185997197,"kyced_addresses":0,"new_unique_addresses":27,"total_burn_fee":3.5843774403434443,"total_eth_volume":9.101505111598318,"total_trade":60,"total_usd_amount":4738.284168671162,"unique_addresses":40,"usd_per_trade":78.97140281118602},"1522022400000":{"eth_per_trade":0.1305336778977258,"kyced_addresses":0,"new_unique_addresses":13,"total_burn_fee":1.2758795269915402,"total_eth_volume":2.3496062021590642,"total_trade":18,"total_usd_amount":1230.3892752776494,"unique_addresses":18,"usd_per_trade":68.35495973764719}},"success":true}
```

### Get a list of wallet that has ever traded with core
```
<host>:8000/get-wallet-address 
GET request

URL params:
  Nil
```


```
curl -x GET http://localhost:8000/get-wallet-address
```
response
```
{"data":["0xb9e29984fe50602e7a619662ebed4f90d93824c7","0xf1aa99c69715f423086008eb9d06dc1e35cc504d"],"success":true}
``` 

### Get exchanges status
```
<host>:8000/get-exchange-status
GET request
```

eg:
```
curl -x GET http://localhost:8000/get-exchange-status
```

response:
```
{"data":{"binance":{"timestamp":1521532176702,"status":true},"bittrex":{"timestamp":1521532176704,"status":true},"huobi":{"timestamp":1521532176703,"status":true}},"success":true}
```

### Update exchanges status
```
<host>:8000/update-exchange-status
POST request

params: 
exchange (string): exchange name (eg: 'binance')
status (bool): true (up), false (down)
timestamp (integer): timestamp of the exchange status
```

eg:
```
curl -X POST \
  http://localhost:8000/update-exchange-status \
  -H 'content-type: multipart/form-data' \
  -F exchange=binance \
  -F status=false
```

### Get country stats
```
<host>:8000/get-country-stats
GET request
params:
 - fromTime (integer) - from timestamp (millisecond)
 - toTime (integer) - to timestamp (millisecond)
 - country (string) - internatinal country 
 - timezone (integer) - timezone to get country stats from -11 to 14
```

response: 
```
{"data":{"1522368000000":{"eth_per_trade":1.1759348083481784,"kyced_addresses":0,"new_unique_addresses":23,"total_burn_fee":40.10625390027786,"total_eth_volume":51.741131567319854,"total_trade":44,"total_usd_amount":19804.392524011764,"unique_addresses":26,"usd_per_trade":450.09983009117644}},"success":true}
```

### Get heatmap - list of countries sort by total ETH value
```
<host>:8000/get-heat-map
GET request
params:
 - fromTime (integer) - from timestamp (millisecond)
 - toTime (integer) - to timestamp (millisecond)
 - timezone (integer) - timezone to get country stats from -11 to 14
```

response:
```
{"data":[{"country":"US","total_eth_value":51.741131567319854,"total_fiat_value":19804.392524011764},{"country":"unknown","total_eth_value":31.28130484378119,"total_fiat_value":12268.937507634406},{"country":"TW","total_eth_value":15,"total_fiat_value":5916.6900000000005},{"country":"KR","total_eth_value":13.280037553077175,"total_fiat_value":5016.70456645198},{"country":"JP","total_eth_value":10.277090646,"total_fiat_value":3857.271305900826},{"country":"TH","total_eth_value":8.241091466923997,"total_fiat_value":3195.368602817533},{"country":"CA","total_eth_value":3.8122812821017558,"total_fiat_value":1445.8819158742285},{"country":"AU","total_eth_value":2.6,"total_fiat_value":969.02},{"country":"DE","total_eth_value":1.823287,"total_fiat_value":697.502009413},{"country":"ID","total_eth_value":1.7178731840736186,"total_fiat_value":674.8439050493492},{"country":"RO","total_eth_value":1.4009999999999998,"total_fiat_value":529.075415},{"country":"VN","total_eth_value":1.3951777988339262,"total_fiat_value":548.8376078547749},{"country":"CN","total_eth_value":1.0121575386522288,"total_fiat_value":401.6824093511598},{"country":"PL","total_eth_value":0.379699,"total_fiat_value":144.141714079},{"country":"FR","total_eth_value":0.319624,"total_fiat_value":122.92586391999998},{"country":"SG","total_eth_value":0.15642985716526572,"total_fiat_value":64.06928945889221},{"country":"ES","total_eth_value":0.09344946,"total_fiat_value":35.176806429959996},{"country":"XX","total_eth_value":0.09,"total_fiat_value":36.86148},{"country":"IN","total_eth_value":0.0714026952146661,"total_fiat_value":27.977050948875906},{"country":"AR","total_eth_value":0.02751473,"total_fiat_value":10.92519129691},{"country":"RU","total_eth_value":0.024162,"total_fiat_value":9.61210186},{"country":"SE","total_eth_value":0.023,"total_fiat_value":9.132541},{"country":"LV","total_eth_value":0.01,"total_fiat_value":3.9209899999999998},{"country":"AL","total_eth_value":0.003,"total_fiat_value":1.126449}],"success":true}
```

### Update Price Analytic Data - (signing required) set a record marking the condition because of which the set price is called. 
```
<host>:8000/update-price-analytic-data
POST request
params:
 - timestamp - the timestamp of the action (real time ) in millisecond
 - value - the json enconded object to save. 

Note: the data sent over must be encoded in Json in order to make it valid for output operation
  In Python, the data would be encoded as:
   data = {"timestamp": timestamp, "value": json.dumps(analytic_data)} 
 ```

response:
```
on success:
{"success":true}
on failure:
{"success":false,
 "reason":<error>}
```

### Get Price Analytic Data - (signing required) list of price analytic data, sorted by timestamp 
```
<host>:8000/get-get-price-analytic-data
GET request
params:
 - fromTime (integer) - from timestamp (millisecond)
 - toTime (integer) - to timestamp (millisecond)
```
example:
```
curl -x GET \
  http://localhost:8000/get-price-analytic-data?fromTime=1522753160000&toTime=1522755792000
```
 
response:
```
{
  "data": [
    {
      "Timestamp": 1522755271000,
      "Data": {
        "block_expiration": false,
        "trigger_price_update": true,
        "triggering_tokens_list": [
          {
            "ask_price": 0.002,
            "bid_price": 0.003,
            "mid afp_old_price": 0.34555,
            "mid_afp_price": 0.6555,
            "min_spread": 0.233,
            "token": "OMG"
          },
          {
            "ask_price": 0.004,
            "bid_price": 0.005,
            "mid afp_old_price": 0.21555,
            "mid_afp_price": 0.4355,
            "min_spread": 0.133,
            "token": "KNC"
          }
        ]
      }
    }
  ],
  "success": true
}
```

### Update exchange notifications 
```
<host>:8000/exchange-notification
POST request
params:
 - exchange (string) - exchange name
 - action (string) - action name
 - token (string) - token pair
 - fromTime (integer) - from timestamp
 - toTime (integer) - to timestamp
 - isWarning (bool) - is exchange warning or not
 - msg (string) - message for the notification
```

response:
```
  {
    "success": true
  }
```

### Get exchange notifications
```
<host>:8000/exchange-notifications
GET request
```

response:
```
{"data":{"binance":{"trade":{"OMG":{"fromTime":123,"toTime":125,"isWarning":true,"msg":"3 times"}}}},"success":true}
```

### Get reserve volume
```
<host>:8000/exchange-notifications
GET request
URL Params:
  - fromTime (integer): millisecond
  - toTime (integer): millisecond
  - token (string): name of token to get volume (eg: ETH)
  - reserveAddr (string): reserve address to get volume of token
  - freq (string): frequency to get volume ("M", "H", "D" - Minute, Hour, Day)
```

example:
```
curl -x GET \
http://localhost:8000/get-reserve-volume?fromTime=1522540800000&toTime=1522627200000&freq=D&token=KNC&reserveAddr=0x63825c174ab367968EC60f061753D3bbD36A0D8F
```

response:
```
{"data":{"1522540800000":{"eth_amount":9.971150530912206,"usd_amount":3838.6105908493496,"volume":3945.5899585215247},"1522627200000":{"eth_amount":14.749439804645423,"usd_amount":5766.650333669346,"volume":5884.90733954939}},"success":true}
```

### set stable token params - (signing required)
```
<host>:8000/set-stable-token-params
POST request
URL Params:
  - value (string) : the json enconded string, represent a map (string : interface)
```


response:
```
on success:
{"success":true}
on failure:
{"success":false,
 "reason":<error>}
```
### confirm stable token params - (signing required)
```
<host>:8000/confirm-stable-token-params
POST request
URL Params:
  - value (string) : the json enconded string, represent a map (string : interface), must be equal to current pending.
```


response:
```
on success:
{"success":true}
on failure:
{"success":false,
 "reason":<error>}
```

### reject stable token params - (signing required)
```
<host>:8000/reject-stable-token-params
POST request
URL Params:
  nil
```


response:
```
on success:
{"success":true}
on failure:
{"success":false,
 "reason":<error>}
```
### Get pending stable token params- (signing required) return the current pending stable token params
```
<host>:8000/pending-stable-token-params
GET request
params:
  - nonce (uint64) : the nonce to conform to signing requirement
```
example:
```
curl -x GET \
  http://localhost:8000/pending-token-params?nonce=111111
```
 
response:
```
{
  "data": {
    "DGX": {
      "AskSpread": 50,
      "BidSpread": 50,
      "PriceUpdateThreshold": 0.1
    }
  },
  "success": true
}
```

### Get stable token params- (signing required) return the current confirmed stable token params
```
<host>:8000/stable-token-params
GET request
params:
  - nonce (uint64) : the nonce to conform to signing requirement
```
example:
```
curl -x GET \
  http://localhost:8000/stable-token-params?nonce=111111
```
 
response:
```
{
  "data": {
    "DGX": {
      "AskSpread": 50,
      "BidSpread": 50,
      "PriceUpdateThreshold": 0.1
    }
  },
  "success": true
}
```
### Get heat map for token
```
<host>:8000/get-token-heatmap
URL Params:
  - fromTime (integer): millisecond
  - toTime (integer): millisecond
  - token (string): name of token to get heatmap
  - freq (string): frequencty to get volume ("M", "H", "D" - Minute, Hour, Day)
```

example:
```
curl -x GET \
http://localhost:8000/get-token-heatmap?fromTime=1518307200000&token=EOS&freq=D&toTime=1518911999999
```

response:
```
{"data":[{"country":"US","volume":2883.620428022146,"eth_volume":29.97000000311978,"usd_volume":28584.013502715607},{"country":"unknown","volume":663.7763113279779,"eth_volume":6.848675774186141,"usd_volume":5710.033060275751},{"country":"JP","volume":189.38349888667832,"eth_volume":1.99,"usd_volume":1881.86987},{"country":"KR","volume":93.83012247596538,"eth_volume":1,"usd_volume":857.766},{"country":"SI","volume":73.000042,"eth_volume":0.7584920000216375,"usd_volume":696.7810908998771},{"country":"IL","volume":9.757144977962138,"eth_volume":0.1,"usd_volume":85.47670000000001},{"country":"TH","volume":9.459436814264475,"eth_volume":0.1,"usd_volume":84.1759},{"country":"DE","volume":9.311558446913438,"eth_volume":0.09904,"usd_volume":85.93066944},{"country":"VN","volume":1.8918873628528947,"eth_volume":0.019789900740301923,"usd_volume":16.536080320374314}],"success":true}
```

### Get gold data
```
<host>:8000/gold-feed
```
response:
```
{"data":{"Timestamp":1526923808631,"DGX":{"Valid":true,"Timestamp":0,"success":"","data":[{"symbol":"DGXETH","price":0.06676463,"time":1526923801},{"symbol":"ETHUSD","price":694.4,"time":1526923801},{"symbol":"ETHSGD","price":931.89,"time":1526923801},{"symbol":"DGXUSD","price":46.36,"time":1526923801},{"symbol":"EURUSD","price":1.17732,"time":1526923801},{"symbol":"USDSGD","price":1.34201,"time":1526923801},{"symbol":"XAUUSD","price":1291.468,"time":1526923801},{"symbol":"USDJPY","price":111.061,"time":1526923801}],"Error":""},"OneForgeETH":{"Value":1.85646,"Text":"1 XAU is worth 1.85646 ETH","Timestamp":1526923803,"Error":false,"Message":""},"OneForgeUSD":{"Value":1291.57,"Text":"1 XAU is worth 1291.57 USD","Timestamp":1526923803,"Error":false,"Message":""},"GDAX":{"Valid":true,"Error":"","trade_id":34527604,"price":"695.56000000","size":"0.00894700","bid":"695.55","ask":"695.56","volume":"50497.82498957","time":"2018-05-21T17:30:04.729000Z"},"Kraken":{"Valid":true,"network_error":"","error":[],"result":{"XETHZUSD":{"a":["696.66000","1","1.000"],"b":["696.33000","4","4.000"],"c":["696.33000","0.10776064"],"v":["13536.83019524","16999.30348103"],"p":["707.93621","710.18316"],"t":[5361,8276],"l":["693.97000","693.97000"],"h":["721.38000","724.80000"],"o":"715.65000"}}},"Gemini":{"Valid":true,"Error":"","bid":"694.50","ask":"695.55","volume":{"ETH":"11418.5646926","USD":"8064891.13775284649999999999999999999704534","timestamp":1526923800000},"last":"695.36"}},"success":true}
```

### Get BTC data

```
<host>:8000/btc-feed
response:
```json
{
  "data": {
    "Timestamp": 1541571292437,
    "bitfinex": {
      "Valid": true,
      "Error": "",
      "mid": "0.0332995",
      "bid": "0.03329",
      "ask": "0.033309",
      "last_price": "0.033299",
      "low": "0.032559",
      "high": "0.034036",
      "volume": "29077.44919025",
      "timestamp": "1541571292.1917806"
    },
    "binance": {
      "Valid": true,
      "Error": "",
      "symbol": "ETHBTC",
      "bidPrice": "0.03328500",
      "bidQty": "2.04600000",
      "askPrice": "0.03328700",
      "askQty": "0.17800000"
    }
  },
  "success": true
}
```

### Set Feed Configuration - (signing required)

```
<host>:8000/set-feed-configuration
URL Params:
- data: JSON encoded string {"name": <feed_name>, "enabled": true/false}
```

response:

```json
{"success":true}
```

### Get Feed Configuration - (signing required)


```
<host>:8000/get-feed-configuration
```

response:

```json
{
  "data": [
    {
      "name": "DGX",
      "enabled": true
    },
    {
      "name": "OneForgeETH",
      "enabled": true
    },
    {
      "name": "OneForgeUSD",
      "enabled": true
    },
    {
      "name": "GDAX",
      "enabled": true
    },
    {
      "name": "Kraken",
      "enabled": true
    },
    {
      "name": "Gemini",
      "enabled": true
    },
    {
      "name": "bitfinex",
      "enabled": true
    },
    {
      "name": "binance",
      "enabled": true
    }
  ],
  "success": true
}
```

### set target quantity v2 - (signing required)
```
<host>:8000/v2/settargetqty
POST request
URL Params:
  - value (string) : the json enconded string, represent a map (string : interface)
```


response:
```
on success:
{"success":true}
on failure:
{"success":false,
 "reason":<error>}
```
### confirm target quantity v2- (signing required)
```
<host>:8000/v2/confirmtargetqty
POST request
URL Params:
  - value (string) : the json enconded string, represent a map (string : interface), must be equal to current pending.
```


response:
```
on success:
{"success":true}
on failure:
{"success":false,
 "reason":<error>}
```

### cancel set target quantity v2- (signing required)
```
<host>:8000/v2/canceltargetqty
POST request
URL Params:
  nil
```


response:
```
on success:
{"success":true}
on failure:
{"success":false,
 "reason":<error>}
```
### Get pending target quantity - (signing required) return the current pending target quantity 
```
<host>:8000/v2/pendingtargetqty
GET request
params:
  - nonce (uint64) : the nonce to conform to signing requirement
```
example:
```
curl -x GET \
  http://localhost:8000/v2/pendingtargetqty?nonce=111111
```
 
response:
```
{
  "data": {
     "OMG" : {
        "TotalTarget": 1500,
        "ReserveTarget": 1005,
        "RebalanceThreshold": 0.33,
        "TransferThreshold": 0.2
    }
  },
  "success": true
}
```

### Get target quantity - (signing required) return the current confirmed target quantity 
```
<host>:8000/v2/targetqty
GET request
params:
  - nonce (uint64) : the nonce to conform to signing requirement
```
example:
```
curl -x GET \
  http://localhost:8000/v2/targetqty?nonce=111111
```
 
response:
```
{
  "data": {
    "OMG" : {
      "TotalTarget": 1500,
      "ReserveTarget": 1005,
      "RebalanceThreshold": 0.33,
      "TransferThreshold": 0.2
    }
  },
  "success": true
}
```

### Get pwis equation v2 - (signing required)

```
<host>:8000/v2/pwis-equation
GET request
```

eg:

```
curl -X "GET" "http://localhost:8000/v2/pwis-equation" \
     -H 'Content-Type: application/x-www-form-urlencoded' \
```

response:
```
{
  "data": {
    "EOS": {
      "ask": {
        "a": 800,
        "b": 600,
        "c": 0,
        "min_min_spread": 0,
        "price_multiply_factor": 0
      },
      "bid": {
        "a": 750,
        "b": 500,
        "c": 0,
        "min_min_spread": 0,
        "price_multiply_factor": 0
      }
    },
    "ETH": {
      "ask": {
        "a": 800,
        "b": 600,
        "c": 0,
        "min_min_spread": 0,
        "price_multiply_factor": 0
      },
      "bid": {
        "a": 750,
        "b": 500,
        "c": 0,
        "min_min_spread": 0,
        "price_multiply_factor": 0
      }
    }
  },
  "success": true
}
```

### Get pending pwis equation v2 - (signing required)

```
<host>:8000/v2/pending-pwis-equation
GET request

```

eg:
```
curl -X "GET" "http://localhost:8000/v2/pending-pwis-equation" \
     -H 'Content-Type: application/x-www-form-urlencoded' \
```

response:
```
{
  "data": {
    "EOS": {
      "ask": {
        "a": 800,
        "b": 600,
        "c": 0,
        "min_min_spread": 0,
        "price_multiply_factor": 0
      },
      "bid": {
        "a": 750,
        "b": 500,
        "c": 0,
        "min_min_spread": 0,
        "price_multiply_factor": 0
      }
    },
    "ETH": {
      "ask": {
        "a": 800,
        "b": 600,
        "c": 0,
        "min_min_spread": 0,
        "price_multiply_factor": 0
      },
      "bid": {
        "a": 750,
        "b": 500,
        "c": 0,
        "min_min_spread": 0,
        "price_multiply_factor": 0
      }
    }
  },
  "success": true
}
```

### Set pwis equation v2 - (signing required)

```
<host>:8000/v2/set-pwis-equation
POST request
Post form: json encoding data of pwis equation
```

eg:

```
curl -X "POST" "http://localhost:8000/v2/set-pwis-equation" \
     -H 'Content-Type: application/x-www-form-urlencoded' \
     --data-urlencode "data={
  \"EOS\": {
    \"bid\": {
      \"a\": 750,
      \"b\": 500,
      \"c\": 0,
      \"min_min_spread\": 0,
      \"price_multiply_factor\": 0
    },
    \"ask\": {
      \"a\": 800,
      \"b\": 600,
      \"c\": 0,
      \"min_min_spread\": 0,
      \"price_multiply_factor\": 0
    }
  },
  \"ETH\": {
    \"bid\": {
      \"a\": 750,
      \"b\": 500,
      \"c\": 0,
      \"min_min_spread\": 0,
      \"price_multiply_factor\": 0
    },
    \"ask\": {
      \"a\": 800,
      \"b\": 600,
      \"c\": 0,
      \"min_min_spread\": 0,
      \"price_multiply_factor\": 0
    }
  }
}"
```

response

```
  {
    "success": true,
  }
```

### Confirm pending pwis equation v2 - (signing required)

```
<host>:8000/v2/confirm-pwis-equation
POST request
Post form: json encoding data of pwis equation
```

eg

```
curl -X "POST" "http://localhost:8000/v2/confirm-pwis-equation" \
     -H 'Content-Type: application/x-www-form-urlencoded' \
     --data-urlencode "data={
  \"EOS\": {
    \"bid\": {
      \"a\": 750,
      \"b\": 500,
      \"c\": 0,
      \"min_min_spread\": 0,
      \"price_multiply_factor\": 0
    },
    \"ask\": {
      \"a\": 800,
      \"b\": 600,
      \"c\": 0,
      \"min_min_spread\": 0,
      \"price_multiply_factor\": 0
    }
  },
  \"ETH\": {
    \"bid\": {
      \"a\": 750,
      \"b\": 500,
      \"c\": 0,
      \"min_min_spread\": 0,
      \"price_multiply_factor\": 0
    },
    \"ask\": {
      \"a\": 800,
      \"b\": 600,
      \"c\": 0,
      \"min_min_spread\": 0,
      \"price_multiply_factor\": 0
    }
  }
}"
```

response

```
  {
    "success": true,
  }
```

### Reject pending pwis equation v2 - (signing required)

```
<host>:8000/v2/reject-pwis-equation
POST request
```

eg

```
curl -X "POST" "http://localhost:8000/v2/reject-pwis-equation" \
     -H 'Content-Type: application/x-www-form-urlencoded'
```

response

```
  {
    "success": true,
  }
```

### Get rebalance quadratic - (signing required)

```
<host>:8000/rebalance-quadratic
GET request

```

eg:
```
curl -X "GET" "http://localhost:8000/rebalance-quadratic" \
     -H 'Content-Type: application/x-www-form-urlencoded' \
```

response:
```
{
  "data": {
    "EOS": {
      "rebalance_quadratic": {
        "a": 800,
        "b": 600,
        "c": 0
      }
    },
    "ETH": {
      "rebalance_quadratic": {
        "a": 750,
        "b": 500,
        "c": 0
      }
    }
  },
  "success": true
}
```

### Set rebalance quadratic equation - (signing required)

```
<host>:8000/set-rebalance-quadratic
POST request
Post form: json encoding data of rebalance quadratic equation
```

eg:

```
curl -X "POST" "http://localhost:8000/set-rebalance-quadratic" \
     -H 'Content-Type: application/x-www-form-urlencoded' \
     --data-urlencode "data={
  "EOS":{
    "rebalance_quadratic": {
      "a": 750,
      "b": 500,
      "c": 0,
    }
  },
  "ETH": {
    "rebalance_quadratic": {
      "a": 750,
      "b": 500,
      "c": 0,
    }
  }
}"
```

response

```

### Get pending rebalance quadratic - (signing required)

```
<host>:8000/pending-rebalance-quadratic
GET request

```

eg:
```
curl -X "GET" "http://localhost:8000/pending-rebalance-quadratic" \
     -H 'Content-Type: application/x-www-form-urlencoded' \
```

response:
```
{
  "data": {
    "EOS": {
      "rebalance_quadratic": {
        "a": 800,
        "b": 600,
        "c": 0
      }
    },
    "ETH": {
      "rebalance_quadratic": {
        "a": 750,
        "b": 500,
        "c": 0
      }
    }
  },
  "success": true
}
```


### Confirm rebalance quadratic equation - (signing required)

```
<host>:8000/confirm-rebalance-quadratic
POST request
Post form: json encoding data of pwis equation
```

eg

```
curl -X "POST" "http://localhost:8000/confirm-rebalance-quadratic" \
     -H 'Content-Type: application/x-www-form-urlencoded' \
     --data-urlencode "data={
  "EOS":{
    "rebalance_quadratic": {
      "a": 750,
      "b": 500,
      "c": 0,
    }
  },
  "ETH": {
    "rebalance_quadratic": {
      "a": 750,
      "b": 500,
      "c": 0,
    }
  }
}"
```

response

```
  {
    "success": true,
  }
```

### Reject rebalance quadrtic equation - (signing required)

```
<host>:8000/reject-rebalance-quadratic
POST request
```

eg

```
curl -X "POST" "http://localhost:8000/reject-rebalance-quadratic" \
     -H 'Content-Type: application/x-www-form-urlencoded'
```

response

```
  {
    "success": true,
  }
```


### Setting APIs
#### Token related APIs

##### Set token update - (signing required) Prepare token update and store the request as pending
POST request 
Post form: {"data" : "JSON enconding of token update Object"}
```
<host>:8000/setting/set-token-update
```
**Note**: 
- The API allow user to update token settings and its status. Hence can be used both for **list** and **delist** a token, as well as 
to do minor modification for the token setting. 
To list a token, it active status is set to true. To delist a token, both its internal and active status is set to false.
- This data is in the form of a map tokenID:tokenUpdate which allows mutiple token updates at once
- It also allows mutiple requests, for example, one request update OMG, the other update KNC. Both these 
requests will be aggregate in to a list of token to be listed. These can be overwritten as well : if there 
are two requests update KNC, the later will overwite the ealier.  
- If a token is marked as internal, it will be required to come with exchange setting( fee, min deposit, 
exchange precision limit, deposit address) , and metric settings (pwis, targetQty). Since rebalance quadratic
data can be zero value, it is optional. 
- If exchange precision limit (tokenUpdate.Exchange.Info) is null, It can be queried from exchange and 
set automatically for the pair (token-ETH). If this data is available in the request,
it will be prioritize over the exchange queried data.
- In addition, if the update contain any Internal token, that token must be available in Smart contract
in order to update its indices. 
- The tokenID from the map object will overwrite the token object's ID. Hence this token object ID inside the request is optional.

Example: This request will list token OMG and NEO. OMG is internal, NEO is external. 

``` 
curl -X "POST" "http://localhost:8000/setting/set-token-update" \
     -H 'Content-Type: application/x-www-form-urlencoded'\
     --data-urlencode "data={  
      \"OMG\": {
        \"token\": {
          \"id\": \"OMG\",
          \"name\": \"OmisexGO\",
          \"decimals\": 18,
          \"address\": \"0xd26114cd6EE289AccF82350c8d8487fedB8A0C07\",
          \"internal\": true,
          \"active\": true
        },
        \"exchanges\": {
          \"binance\": {
            \"deposit_address\": \"0x22222222222222222222222222222222222\",
            \"fee\": {
              \"withdraw\": 0.2,
              \"deposit\": 0.3
            },
            \"min_deposit\": 4
          }
        },
        \"pwis_equation\": {
          \"ask\": {
            \"a\": 800,
            \"b\": 600,
            \"c\": 0,
            \"min_min_spread\": 0,
            \"price_multiply_factor\": 0
          },
          \"bid\": {
            \"a\": 750,
            \"b\": 500,
            \"c\": 0,
            \"min_min_spread\": 0,
            \"price_multiply_factor\": 0
          }
        },
        \"target_qty\": {
          \"set_target\": {
            \"total_target\": 0,
            \"reserve_target\": 0,
            \"rebalance_threshold\": 0,
            \"transfer_threshold\": 0
          }
        },
        \"rebalance_quadratic\": {
          \"rebalance_quadratic\": {
            \"a\": 1,
            \"b\": 2,
            \"c\": 3
          }
        }
      },
      \"NEO\": {
        \"Token\": {
          \"id\": \"NEO\",
          \"name\": \"Request\",
          \"decimals\": 18,
          \"address\": \"0x8f8221afbb33998d8584a2b05749ba73c37a938a\",
          \"internal\": false,
          \"active\": true
        }
      }
    }"

```
response

```
on success:
{"success":true}
on failure:
{"success":false,
 "reason":<error>}
```

##### Get pending token update - (singing required) Return the current pending token updates information
GET request

``` 
<host>:8000/setting/pending-token-update
```

Example
``` 
curl -X "GET" "http://localhost:8000/setting/pending-token-update"
```

response 
```
{
  "data": {
    "NEO": {
      "token": {
        "id": "NEO",
        "name": "Request",
        "address": "0x8f8221afbb33998d8584a2b05749ba73c37a938a",
        "decimals": 18,
        "active": true,
        "internal": false,
        "last_activation_change": 0
      },
      "exchanges": null,
      "pwis_equation": null,
      "target_qty": {
        "set_target": {
          "total_target": 0,
          "reserve_target": 0,
          "rebalance_threshold": 0,
          "transfer_threshold": 0
        }
      },
      "rebalance_quadratic": {
        "rebalance_quadratic": {
          "a": 0,
          "b": 0,
          "c": 0
        }
      }
    },
    "OMG": {
      "token": {
        "id": "OMG",
        "name": "OmisexGO",
        "address": "0xd26114cd6EE289AccF82350c8d8487fedB8A0C07",
        "decimals": 18,
        "active": true,
        "internal": true,
        "last_activation_change": 0
      },
      "exchanges": {
        "binance": {
          "deposit_address": "",
          "exchange_info": {
            "OMG-ETH": {
              "precision": {
                "amount": 2,
                "price": 6
              },
              "amount_limit": {
                "min": 0.01,
                "max": 90000000
              },
              "price_limit": {
                "min": 0.001611,
                "max": 0.16103
              },
              "min_notional": 0.01
            }
          },
          "fee": {
            "withdraw": 0.2,
            "deposit": 0.3
          },
          "min_deposit": 0
        }
      },
      "pwis_equation": {
        "ask": {
          "a": 800,
          "b": 600,
          "c": 0,
          "min_min_spread": 0,
          "price_multiply_factor": 0
        },
        "bid": {
          "a": 750,
          "b": 500,
          "c": 0,
          "min_min_spread": 0,
          "price_multiply_factor": 0
        }
      },
      "target_qty": {
        "set_target": {
          "total_target": 1,
          "reserve_target": 2,
          "rebalance_threshold": 0,
          "transfer_threshold": 0
        }
      },
      "rebalance_quadratic": {
        "rebalance_quadratic": {
          "a": 1,
          "b": 2,
          "c": 3
        }

##### Confirm token update - (signing required) Confirm token update and apply all the change to core.
POST request 
Post form: {"data" : "JSON enconding of token update Object"}
Note: This data is similar to token update, but all field must be the same as the current pending. 
```
<host>:8000/setting/confirm-token-update
```

Example 

``` 
curl -X "POST" "http://localhost:8000/setting/confirm-token-update" \
     -H 'Content-Type: application/x-www-form-urlencoded'\
     --data-urlencode "data={    
        \"NEO\": {
          \"token\": {
            \"id\": \"NEO\",
            \"name\": \"Request\",
            \"address\": \"0x8f8221afbb33998d8584a2b05749ba73c37a938a\",
            \"decimals\": 18,
            \"active\": true,
            \"internal\": false
          },
          \"exchanges\": null,
          \"pwis_equation\": null,
          \"target_qty\": {
            \"set_target\": {
              \"total_target\": 0,
              \"reserve_target\": 0,
              \"rebalance_threshold\": 0,
              \"transfer_threshold\": 0
            }
          },
          \"rebalance_quadratic\": {
            \"rebalance_quadratic\": {
              \"a\": 0,
              \"b\": 0,
              \"c\": 0
            }
          }
        },
        \"OMG\": {
          \"token\": {
            \"id\": \"OMG\",
            \"name\": \"OmisexGO\",
            \"address\": \"0xd26114cd6EE289AccF82350c8d8487fedB8A0C07\",
            \"decimals\": 18,
            \"active\": true,
            \"internal\": true
          },
          \"exchanges\": {
            \"binance\": {
              \"deposit_address\": \"0x22222222222222222222222222222222222\",
              \"exchange_info\": {
                \"OMG-ETH\": {
                  \"precision\": {
                    \"amount\": 2,
                    \"price\": 6
                  },
                  \"amount_limit\": {
                    \"min\": 0.01,
                    \"max\": 90000000
                  },
                  \"price_limit\": {
                    \"min\": 0.000001,
                    \"max\": 100000
                  },
                  \"min_notional\": 0.01
                }
              },
              \"fee\": {
                \"withdraw\": 0.2,
                \"deposit\": 0.3
              },
              \"min_deposit\": 4
            }
          },
          \"pwis_equation\": {
            \"ask\": {
              \"a\": 800,
              \"b\": 600,
              \"c\": 0,
              \"min_min_spread\": 0,
              \"price_multiply_factor\": 0
            },
            \"bid\": {
              \"a\": 750,
              \"b\": 500,
              \"c\": 0,
              \"min_min_spread\": 0,
              \"price_multiply_factor\": 0
            }
          },
          \"target_qty\": {
            \"set_target\": {
              \"total_target\": 0,
              \"reserve_target\": 0,
              \"rebalance_threshold\": 0,
              \"transfer_threshold\": 0
            }
          },
          \"rebalance_quadratic\": {
            \"rebalance_quadratic\": {
              \"a\": 0,
              \"b\": 0,
              \"c\": 0
            }
          }
        }
    }"

```
response

```
on success:
{"success":true}
on failure:
{"success":false,
 "reason":<error>}
```

##### Reject pending token update - (signing required) reject the update and remove the current pending update
POST request

```
<host>:8000/setting/reject-token-update
```

Example

```
curl -X "POST" "http://localhost:8000/setting/reject-token-update" \
     -H 'Content-Type: application/x-www-form-urlencoded'


on success:
{"success":true}
on failure:
{"success":false,
 "reason":<error>}
```

##### Get Token settings - (signing required) get current token settings of core.
GET request

``` 
<host>:8000/setting/token-settings
```

Example
```
curl -X "GET" "http://localhost:8000/setting/token-settings"
```

response
 
```json
{
  "data": [
    {
      "id": "ABT",
      "name": "",
      "address": "0xb98d4c97425d9908e66e53a6fdf673acca0be986",
      "decimals": 18,
      "active": true,
      "internal": true
    }
  ],
  "success": true
}
```
#### Address related APIs

##### Update address - (signing required) update a single address
POST request 
Post form: {"name" : "Name of the address (reserve, deposit etc...)",
            "address" : "Hex form of the new address"
            "timestamp" (optional) uint64 "this will overwrite version in address setting"  }
Note: This is used to update single address object. For list of address object, use add-address-to-set instead
```
<host>:8000/setting/update-address
```

Example 

```
curl -X "POST" "http://localhost:8000/setting/update-address" \
     -H 'Content-Type: application/x-www-form-urlencoded'\
     --data-urlencode "name=bank"\
     --data-urlencode "address=0x123456789aabbcceeeddff"\
     --data-urlencode "timestamp=1111111111"

```
response

```
on success:
{"success":true}
on failure:
{"success":false,
 "reason":<error>}
```

##### Add address to set- (signing required) add address to a list of address
POST request 
Post form: {"name" : <Name of the address set(oldBurners etc...)>,
            "address" : <Hex form of the new address>
            "timestamp" (optional) uint64 <this will overwrite version in address setting> }
```
<host>:8000/setting/add-address-to-set
```

Example 

```
curl -X "POST" "http://localhost:8000/setting/add-address-to-set" \
     -H 'Content-Type: application/x-www-form-urlencoded'\
     --data-urlencode "name="third_party_reserves"\
     --data-urlencode "address=0x123456789aabbcceeeddff" 

```
response

```
on success:
{"success":true}
on failure:
{"success":false,
 "reason":<error>}
```

#### Exchange related APIs

##### Update exchange fee - (signing required) update one exchange fee setting
POST request 
Post form: {"name" : <Name of the exchange (binance, huobi etc...)>,
            "data" : <JSON encoded form of fee setting >
            "timestamp" (optional) uint64 <this will overwrite version in exchange setting> }
}
**Note**: 
UpdateFee will merge the new fee setting to the current fee setting,
Any different key will be overwriten from new fee to current fee. This allows update
one single token's exchange fee on a destined exchange.
UpdateFee will not be mutiplied by any value, so please prepare a big enough number to avoid exchange's fee increasing.
```
<host>:8000/setting/update-exchange-fee
```

Example 

```
  curl -X "POST" "http://localhost:8000/setting/update-exchange-fee" \
     -H 'Content-Type: application/x-www-form-urlencoded'\
     --data-urlencode "name=binance"\
     --data-urlencode "data= {
      \"Trading\": {
        \"maker\": 0.001,
        \"taker\": 0.001
      },
      \"Funding\": {
        \"Withdraw\": {
          \"ZEC\": 0.005,
          \"ZIL\": 100,
          \"ZRX\": 5.8
        },
        \"Deposit\": {
          \"ZEC\": 0,
          \"ZIL\": 0,
          \"ZRX\": 2
        }
      }
    }"
```
##### Update exchange mindeposit - (signing required) update one exchange min deposit
POST request 
Post form: {"name" : <Name of the exchange (binance, huobi etc...)>,
            "data" : <JSON encoded form of min deposit>
            "timestamp" (optional) uint64 <this will overwrite version in exchange setting> }

**Note**: 
Update Exchange minDeposit will merge the new minDeposit setting to the current minDeposit setting,
Any different key will be overwriten from new minDeposit to current minDeposit. This allows update
one single token's exchange minDeposit on a destined exchange.
minDeposit input will not be mutiplied by any value, so please prepare a big enough number to avoid exchange's minDeposit increasing.
```
<host>:8000/setting/update-exchange-mindeposit
```

Example 

```
  curl -X "POST" "http://localhost:8000/setting/update-exchange-mindeposit" \
     -H 'Content-Type: application/x-www-form-urlencoded'\
     --data-urlencode "name=binance"\
     --data-urlencode "data= {
      \"POWR\": 0.1,
      \"MANA\": 0.2
    }"
```

#####  Update exchange deposit address - (signing required) update one exchange deposit address
POST request 
Post form: {"name" : <Name of the exchange (binance, huobi etc...)>,
            "data" : <JSON encoded form of a map of token : depositaddress >
            "timestamp" (optional) uint64 <this will overwrite version in exchange setting> }

**Note**: 
Update Exchange deposit address will merge the new deposit address setting to the current deposit address setting,
Any different key will be overwriten from new deposit address to current deposit address. This allows update
one single tokenpair's exchange precision limit on a destined exchange.
```
<host>:8000/setting/update-deposit-address
```

Example 

```shell
  curl -X "POST" "http://localhost:8000/setting/update-deposit-address" \
     -H 'Content-Type: application/x-www-form-urlencoded'\
     --data-urlencode "name=binance"\
     --data-urlencode "data= {
      \"POWR\": \"0x778599Dd7893C8166D313F0F9B5F6cbF7536c293\"
    }"
```

#####  Update exchange info - (signing required) update one exchange's info

POST request 
Post form: {"name" : <Name of the exchange (binance, huobi etc...)>,
            "data" : <JSON encoded form of exchange info >
            "timestamp" (optional) uint64 <this will overwrite version in exchange setting> }
}
**Note**: 
Update Exchange minDeposit will merge the new exchange info setting to the current exchange info setting,
Any different key will be overwriten from new exchange info to current exchange info. This allows update
one single token's exchange minDeposit on a destined exchange.

```shell
<host>:8000/setting/update-exchange-info
```

Example 

```shell
  curl -X "POST" "http://localhost:8000/setting/update-exchange-info" \
     -H 'Content-Type: application/x-www-form-urlencoded'\
     --data-urlencode "name=binance"\
     --data-urlencode "data= {
      \"LINK-ETH\": {
        \"precision\": {
          \"amount\": 0,
          \"price\": 8
        },
        \"amount_limit\": {
          \"min\": 1,
          \"max\": 90000000
        },
        \"price_limit\": {
          \"min\": 1e-8,
          \"max\": 120000
        },
        \"min_notional\": 0.01
      }
    }"
```

##### Get all settings - (signing required) return all current running setting of core
GET request

```shell
<host>:8000/setting/all-settings
```

Example

```shell
curl -X "GET" "http://localhost:8000/setting/all-settings"
```

Response

```json
{
  "data": {
    "Addresses": {
      "Addresses": {
        "bank": "",
        "burner": "0xed4f53268bfdff39b36e8786247ba3a02cf34b04",
        "deposit_operator": "0xEDd15B61505180B3A0C25B193dF27eF10214D851",
        "intermediate_operator": "0x13922F1857C0677F79e4BbB16Ad2c49fAa620829",
        "internal_network": "0x91a502c678605fbce581eae053319747482276b9",
        "network": "0x818e6fecd516ecc3849daf6845e3ec868087b755",
        "old_burners": [
          "0x07f6e905f2a1559cd9fd43cb92f8a1062a3ca706",
          "0x4e89bc8484b2c454f2f7b25b612b648c45e14a8e"
        ],
        "old_networks": [
          "0x964f35fae36d75b1e72770e244f6595b68508cf5"
        ],
        "pricing": "0x798abda6cc246d0edba912092a2a3dbd3d11191b",
        "pricing_operator": "0x760d30979EB313A2d23C53E4Fb55986183B0ffd9",
        "reserve": "0x63825c174ab367968ec60f061753d3bbd36a0d8f",
        "third_party_reserves": [
          "0x2aab2b157a03915c8a73adae735d0cf51c872f31",
          "0x4d864b5b4f866f65f53cbaad32eb9574760865e6",
          "0x6f50e41885fdc44dbdf7797df0393779a9c0a3a6"
        ],
        "whitelist": "0x6e106a75d369d09a9ea1dcc16da844792aa669a3",
        "wrapper": "0x6172afc8c00c46e0d07ce3af203828198194620a"
      },
      "Version": 1533615419127
    },
    "Tokens": {
      "Tokens": [
        {
          "id": "ABT",
          "name": "ArcBlock",
          "address": "0xb98d4c97425d9908e66e53a6fdf673acca0be986",
          "decimals": 18,
          "active": true,
          "internal": true,
          "last_activation_change": 1533615415641
        },
        {
          "id": "ZIL",
          "name": "Zilliqa",
          "address": "0x05f4a42e251f2d52b8ed15e9fedaacfcef1fad27",
          "decimals": 12,
          "active": true,
          "internal": true,
          "last_activation_change": 1533615415657
        }
      ],
      "Version": 1533615415671
    },
    "Exchanges": {
      "Exchanges": {
        "binance": {
          "deposit_address": {
            "TUSD": "0x44d34a119ba21a42167ff8b77a88f0fc7bb2db90",
            "ZIL": "0xa34c7ac0980c738e4fbf190568f44997a0d4f2dc"
          },
          "min_deposit": {
            "ADA": 0,
            "ADX": 0,
            "AE": 0,
            "ZRX": 0
          },
          "fee": {
            "Trading": {
              "maker": 0.001,
              "taker": 0.001
            },
            "Funding": {
              "Withdraw": {
                "ADA": 2,
                "ADX": 8,
                "AE": 4.6,
                "ZIL": 200,
                "ZRX": 11.6
              },
              "Deposit": {
                "ADA": 0,
                "ADX": 0,
                "AE": 0,
                "AION": 0,
                "ZRX": 0
              }
            }
          },
          "info": {
            "TUSD-ETH": {
              "precision": {
                "amount": 0,
                "price": 8
              },
              "amount_limit": {
                "min": 1,
                "max": 90000000
              },
              "price_limit": {
                "min": 0.0002475,
                "max": 0.0247499
              },
              "min_notional": 0.01
            },
            "ZIL-ETH": {
              "precision": {
                "amount": 0,
                "price": 8
              },
              "amount_limit": {
                "min": 1,
                "max": 90000000
              },
              "price_limit": {
                "min": 0.00001249,
                "max": 0.0012483
              },
              "min_notional": 0.01
            }
          }
        },
        "huobi": {
          "deposit_address": {
            "ABT": "0x0c8fd73eaf6089ef1b91231d0a07d0d2ca2b9d66",
            "WAX": "0x0c8fd73eaf6089ef1b91231d0a07d0d2ca2b9d66"
          },
          "min_deposit": {
            "ABT": 4,
            "APPC": 1,
            "ZIL": 200
          },
          "fee": {
            "Trading": {
              "maker": 0.002,
              "taker": 0.002
            },
            "Funding": {
              "Withdraw": {
                "ZLA": 2,
                "ZRX": 10
              },
              "Deposit": {
                "ZLA": 0,
                "ZRX": 0
              }
            }
          },
          "info": {
            "POLY-ETH": {
              "precision": {
                "amount": 4,
                "price": 6
              },
              "amount_limit": {
                "min": 0,
                "max": 0
              },
              "price_limit": {
                "min": 0,
                "max": 0
              },
              "min_notional": 0.02
            },
            "WAX-ETH": {
              "precision": {
                "amount": 4,
                "price": 6
              },
              "amount_limit": {
                "min": 0,
                "max": 0
              },
              "price_limit": {
                "min": 0,
                "max": 0
              },
              "min_notional": 0.02
            }
          }
        }
      },
      "Version": 1533615419111
    }
  },
  "success": true,
  "timestamp": 1533615425492
}
```

```json
  {
    "success": true
  }
```
