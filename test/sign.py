#python 3.8
import time
import hmac
import hashlib
import base64
import urllib.parse
import requests


timestamp = str(round(time.time() * 1000))
secret = 'SECc3f980869d462e06ac0fef440ab83f633e9aee16531247a9361462e22559e815'
secret_enc = secret.encode('utf-8')
string_to_sign = '{}\n{}'.format(timestamp, secret)
string_to_sign_enc = string_to_sign.encode('utf-8')
hmac_code = hmac.new(secret_enc, string_to_sign_enc, digestmod=hashlib.sha256).digest()
sign = urllib.parse.quote_plus(base64.b64encode(hmac_code))
print(timestamp)
print(sign)


message = """{
    "at": {
        "atMobiles":[
            "15515787648"
        ],
        "isAtAll": false
    },
    "text": {
        "content":"test"
    },
    "msgtype":"text"
}"""
DingUrl = "https://oapi.dingtalk.com/robot/send?access_token=54b3c0269865323ddea151b03e2bc9898dad74ce5239e9364b5f3bad9ec3fb33&timestamp="+ timestamp + "&sign="+sign
result = requests.post(headers={"Content-type":  "application/json;charset=UTF-8"},url=DingUrl,data=message)
result.status_code
print(DingUrl,message)