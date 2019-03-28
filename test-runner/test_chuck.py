import unittest
from chuck import *

class BaseTests(unittest.TestCase):

    def setUp(self):
        self.api = ChuckTester()

    def test_auth_init(self):
        self.api.activate_scenario("scenario-1", "id-1")
        json = self.api.auth_init()
        self.assertEqual(json["code"], "this-is-very-secret-code")

    def test_verify(self):
        self.api.activate_scenario("scenario-1", "id-1")
        json = self.api.verify()
        self.assertEqual(json["code"], "this-is-other-secret-code")

    def test_login(self):
        self.api.activate_scenario("scenario-1", "id-1")
        json = self.api.login()
        self.assertEqual(json["accessToken"], "access-token")
        self.assertEqual(json["refreshToken"], "refresh-token")
        self.assertEqual(json["refreshTokenExpiresIn"], 1209599)

    def test_wrong_url(self):
        self.api.activate_scenario("scenario-1", "id-1")
        code = self.api.wrong_path()
        self.assertEqual(code, 404)

    def test_wrong_query(self):
        self.api.activate_scenario("scenario-1", "id-1")
        code = self.api.wrong_query()
        self.assertEqual(code, 404)

    def test_guest_auth(self):
        self.api.activate_scenario("scenario-2", "id-2")
        code = self.api.guest_auth()
        self.assertEqual(code, 200)

    def test_delete_token(self):
        self.api.activate_scenario("scenario-2", "id-2")
        code = self.api.delete_token()
        self.assertEqual(code, 200)

    def test_auth_init_sc3(self):
        self.api.activate_scenario("scenario-3", "id-3")
        json = self.api.auth_init_sc3()
        self.assertEqual(json["code"], "scenario-3-code")

    def test_preferred_store(self):
        self.api.activate_scenario("scenario-3", "id-3")
        json = self.api.preferred_store()
        self.assertEqual(json["Store"]["Name"], "store-name")



if __name__ == '__main__':
    unittest.main()
