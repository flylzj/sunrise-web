# api

## 获取时间设置

`GET /api/time_interval`

成功后返回

```json
{
    "code": 0,
    "data": {
        "hour": 8,
        "minute": 1
    },
    "message": "ok"
}
```

## 更新时间设置

`POST /api/time_interval`

need

```json
{
    "intervalHour": 8,
    "intervalminute": 1
}
```

成功后返回

```json
{
    "code": 0,
    "message": "ok"
}
```

## 获取邮箱设置

`GET /api/email`

成功后返回

```json
{
    "code": 0,
    "data": {
        "receiver": "",
        "sender": "",
        "sender_pwd": ""
    },
    "message": "ok"
}
```

## 设置邮箱

`POST /api/email`

need

```json
{
    "sender": "1449902124@qq.com",
    "senderPwd": "qwer1234",
    "receiver": "1449902124@qq.com"
}
```

成功后返回

```json
{
    "code": 0,
    "message": "ok"
}
```

## 获取历史信息

`GET /api/history/236081`

成功后返回

```json

{
    "code": 0,
    "data": [
        {
            "ID": 1,
            "Abiid": "236081",
            "Stock": "aa",
            "StockNum": 1,
            "UpdateTime": 1
        }
    ],
    "message": "ok"
}

```

## 获取被监控的商品

`GET /api/good`

成功后返回

```json

{
    "code": 0,
    "data": [
        {
            "ID": 1,
            "Abiid": "236081",
            "MainName": "酷思茗排毒调和茶罐装125克",
            "Subtitle": "KUSMI TEA  KM-DETOX TIN BOX 125G",
            "BrandId": "20920",
            "BrandName": "KUSMI TEA",
            "CategoryId": "0107",
            "CategoryName": "花茶",
            "Price": 186,
            "RealPrice": 204,
            "Stock": "库存:少量",
            "IntStock": 86
        }
    ],
    "message": "ok"
}
```

## 添加需要被监控的商品

`POST /api/good`

need

```json
{
    "abiid": "236081"
}
```

成功后返回

```json

{
    "code": 0,
    "message": "ok"
}
```

## 删除商品

`DELETE /api/good`

need

```json
{
    "abiid": "236081"
}
```

成功后返回

```json
{
    "code": 0,
    "message": "ok"
}
```

## 通过上传excel来批量添加商品

`POST /api/good/upload`

need 

```
file   *.xslx
```

成功后返回

```json
{
    "code": 0,
    "data": [
        [
            "78152",
            "已存在"
        ],
        [
            "121047",
            "已存在"
        ],
        [
            "253865",
            "已存在"
        ],
        [
            "78219",
            "已存在"
        ],
        [
            "77867",
            "已存在"
        ],
        [
            "264003",
            "已存在"
        ],
        [
            "264004",
            "已存在"
        ],
        [
            "118247",
            "已存在"
        ],
        [
            "264269",
            "已存在"
        ],
        [
            "78744",
            "已存在"
        ]
    ],
    "message": "ok"
}
```
