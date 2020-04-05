import random

class Challenge:
    def __init__(self):
        self._url = ''
        self._ans = 'A'
        self.new()

    def __str__(self):
        return '{url},{ans}'.format(url=self._url, ans=self._ans)

    def new(self):
        url='https://vip1.loli.net/2020/04/05/EC7Wnlr2UKoZQvt.png'
        ans='D'

        self._ans = ans
        self._url = url

    def qus(self):
        return self.__str__()

    def ans(self):
        return self._ans

    def url(self):
        return self._url

