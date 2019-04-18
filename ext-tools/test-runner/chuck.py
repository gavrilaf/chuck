import requests


class ChuckTester:
    def __init__(self):
        self.chuck_endpoint = "http://127.0.0.1:8123"
        self.endpoint = "https://test.net"

        self.s = requests.Session()
        self.s.proxies = {'https': 'https://127.0.0.1:8123'}
        self.s.verify = False

    def activate_scenario(self, scenario, id):
        resp = requests.put(self.chuck_endpoint + "/scenario/" + scenario + "/" + id + "/no")

        self.s.headers = {'automation-test-identifier': id}
        # print("Activate scenario {}, result {}".format(scenario, resp))

    # scenario 1
    def auth_init(self):
        resp = self.s.get(self.endpoint + "/v1/authinit?format=json&apikey=1234567&code=7654321")
        if resp.status_code == 200:
            return resp.json()
        else:
            return None

    def verify(self):
        resp = self.s.get(self.endpoint + "/v1/verify?verifier=1a2b3c4d")
        if resp.status_code == 200:
            return resp.json()
        else:
            return None

    def login(self):
        resp = self.s.post(self.endpoint + "/v1/login?format=json&apikey=987654")
        if resp.status_code == 200:
            return resp.json()
        else:
            return None

    # scenario 2
    def guest_auth(self):
        resp = self.s.post(self.endpoint + "/v2/authtoken/guest?format=json&apikey=1234")
        return resp.status_code

    def delete_token(self):
        resp = self.s.delete(self.endpoint + "/v1/authtoken/user-id-111?format=json&apikey=2345")
        return resp.status_code

    # scenario 3
    def auth_init_sc3(self):
        resp = self.s.get(self.endpoint + "/v1/authinit?format=json&apikey=1234567&code=7654321")
        if resp.status_code == 200:
            return resp.json()
        else:
            return None

    def preferred_store(self):
        resp = self.s.get(self.endpoint + "/v1/customer/profile/preferredstore?format=json&apikey=9878")
        if resp.status_code == 200:
            return resp.json()
        else:
            return None

    # invalid requests

    def wrong_path(self):
        resp = self.s.post(self.endpoint + "/v2/login?format=json&apikey=987654")
        return resp.status_code

    def wrong_query(self):
        resp = self.s.post(self.endpoint + "/v1/login?format=xml&apikey=987654")
        return resp.status_code
