import random
import json
from urllib.request import Request, urlopen

class Challenge:
    def __init__(self):
        self._url = ''
        self._ans = 'A'
        self.new()

    def __str__(self):
        return '{url},{ans}'.format(url=self._url, ans=self._ans)



    def new(self):
        with open('config.json', encoding='utf-8') as f:
            config = json.load(f)

        tk = config['tk_url']
        html = urlopen( tk )
        data = html.read()
        strs = str(data)
        strs= strs[2:-1]
        data_json = json.loads(strs)
        url=data_json['Url']
        ans=data_json['Ans']

        self._ans = ans
        self._url = url

    def qus(self):
        return self.__str__()

    def ans(self):
        return self._ans

    def url(self):
        return self._url

