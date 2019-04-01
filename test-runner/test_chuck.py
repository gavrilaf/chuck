import unittest
import threading
import time
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


class MultipleClientsTests(unittest.TestCase):
    def setUp(self):
        self.api = []

    def test_one_scenario_multiple_clients(self):
        for i in range(5):
            self.api.append(ChuckTester())

        self.api[0].activate_scenario("scenario-1", "id-1")
        self.api[1].activate_scenario("scenario-1", "id-2")
        self.api[2].activate_scenario("scenario-1", "id-3")
        self.api[3].activate_scenario("scenario-1", "id-4")
        self.api[4].activate_scenario("scenario-1", "id-5")

        json = self.api[0].auth_init()
        self.assertEqual(json["code"], "this-is-very-secret-code")

        json = self.api[1].verify()
        self.assertEqual(json["code"], "this-is-other-secret-code")

        json = self.api[2].login()
        self.assertEqual(json["accessToken"], "access-token")

        json = self.api[3].auth_init()
        self.assertEqual(json["code"], "this-is-very-secret-code")

        json = self.api[4].verify()
        self.assertEqual(json["code"], "this-is-other-secret-code")

    def test_multiple_scenarios_multiple_clients(self):
        for i in range(6):
            self.api.append(ChuckTester())

        self.api[0].activate_scenario("scenario-1", "id-1")
        self.api[1].activate_scenario("scenario-2", "id-2")
        self.api[2].activate_scenario("scenario-3", "id-3")
        self.api[3].activate_scenario("scenario-1", "id-4")
        self.api[4].activate_scenario("scenario-2", "id-5")
        self.api[5].activate_scenario("scenario-3", "id-6")

        json = self.api[0].auth_init()
        self.assertEqual(json["code"], "this-is-very-secret-code")

        json = self.api[1].verify()
        self.assertIsNone(json)

        code = self.api[1].guest_auth()
        self.assertEqual(code, 200)

        json = self.api[2].auth_init_sc3()
        self.assertEqual(json["code"], "scenario-3-code")

        json = self.api[3].verify()
        self.assertEqual(json["code"], "this-is-other-secret-code")

        code = self.api[4].delete_token()
        self.assertEqual(code, 200)

        json = self.api[5].preferred_store()
        self.assertEqual(json["Store"]["Name"], "store-name")


# Multithread test

class Runner:
    def __init__(self, target, client_id):
        self._target = target
        self._client_id = client_id
        self._result = None

        self._thread = threading.Thread(target = self.run)

    def start(self):
        self._thread.start()

    def run(self):
        self._result = self._target(self._client_id)

    def get_result(self):
        self._thread.join()
        return self._result


def run_scenario_1(client_id):
    api = ChuckTester()
    api.activate_scenario("scenario-1", client_id)

    time.sleep(0)

    json = api.auth_init()
    r1 = json["code"] == "this-is-very-secret-code"

    json = api.verify()
    r2 = json["code"] == "this-is-other-secret-code"

    json = api.login()
    r3 = json["accessToken"] == "access-token"

    return r1, r2, r3


def run_scenario_2(client_id):
    api = ChuckTester()
    api.activate_scenario("scenario-2", client_id)

    time.sleep(0)

    code = api.guest_auth()
    r1 = 200 == code

    code = api.delete_token()
    r2 = 200 == code

    return r1, r2


def run_scenario_3(client_id):
    api = ChuckTester()
    api.activate_scenario("scenario-3", client_id)

    time.sleep(0)

    json = api.auth_init_sc3()
    r1 = json["code"] == "scenario-3-code"

    json = api.preferred_store()
    r2 = json["Store"]["Name"] == "store-name"

    return r1, r2


class MutlithreadTests(unittest.TestCase):

    def test_scenario_1_multiple_clients(self):
        count = 10
        runner = []

        for i in range(count):
            runner.append(Runner(run_scenario_1, "id-{}".format(i)))

        for r in runner:
            r.start()

        for r in runner:
            self.assertEqual(r.get_result(), (True, True, True))

    def test_scenario_2_multiple_clients(self):
        count = 10
        runner = []

        for i in range(count):
            runner.append(Runner(run_scenario_2, "id-{}".format(i)))

        for r in runner:
            r.start()

        for r in runner:
            self.assertEqual(r.get_result(), (True, True))

    def test_scenario_3_multiple_clients(self):
        count = 10
        runner = []

        for i in range(count):
            runner.append(Runner(run_scenario_3, "id-{}".format(i)))

        for r in runner:
            r.start()

        for r in runner:
            self.assertEqual(r.get_result(), (True, True))

    def test_mixed(self):
        def block(client_id):
            r1 = run_scenario_1(client_id)
            r2 = run_scenario_2(client_id)
            r3 = run_scenario_3(client_id)
            return r1 + r2 + r3

        count = 10
        runner = []

        for i in range(count):
            runner.append(Runner(block, "id-{}".format(i)))

        for r in runner:
            r.start()

        for r in runner:
            self.assertEqual(r.get_result(), (True, True, True, True, True, True, True))



if __name__ == '__main__':
    unittest.main()
