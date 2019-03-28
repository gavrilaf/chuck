import unittest
from chuck import *

class BaseTests(unittest.TestCase):

    def setUp(self):
        self.api = ChuckTester()
        self.api.activate_scenario("scenario-1", "id-1")

    def test_auth_init(self):
        json = self.api.auth_init()
        self.assertEqual(json["code"], "this-is-very-secret-code")

    def test_verify(self):
        json = self.api.verify()
        self.assertEqual(json["code"], "this-is-other-secret-code")

    def test_login(self):
        json = self.api.login()
        self.assertEqual(json["accessToken"], "access-token")
        self.assertEqual(json["refreshToken"], "refresh-token")
        self.assertEqual(json["refreshTokenExpiresIn"], 1209599)

    def test_wrong_url(self):
        code = self.api.wrong_path()
        self.assertEqual(code, 404)

    def test_wrong_query(self):
        code = self.api.wrong_query()
        self.assertEqual(code, 404)



if __name__ == '__main__':
    unittest.main()
