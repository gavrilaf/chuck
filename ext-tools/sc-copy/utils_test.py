import unittest
from utils import *


class UtilsTestCase(unittest.TestCase):
    def test_clear_url(self):
        cases = [("https://test.net/guest?format=json&apikey=1234", "https://test.net/guest?format=json&apikey=*"),
                 ("https://test.net/guest?apikey=1234", "https://test.net/guest?apikey=*"),
                 ("https://test.net/guest?apikey=1234&format=json", "https://test.net/guest?apikey=*&format=json"),
                 ("https://test.net/guest?format=json&apikey=1234&code=67-98", "https://test.net/guest?format=json&apikey=*&code=*"),
                 ("https://test.net/guest?format=json&verifier=67-98-1", "https://test.net/guest?format=json&verifier=*"),
                 ("https://test.net/guest?from=0&to=10&from=2019-1-1&to=2019-10-10", "https://test.net/guest?from=0&to=10&from=*&to=*"),
                 ("https://test.net/guest?from=2019-10-10&to=2019-10-12&apikey=1234", "https://test.net/guest?from=*&to=*&apikey=*")]

        for cs in cases:
            cleared = clear_url(cs[0])
            expected = cs[1]

            self.assertEqual(expected, cleared)


if __name__ == '__main__':
    unittest.main()
