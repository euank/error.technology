import httplib
import traceback
import os,sys

def webscale_errorhook(excType, excValue, traceback):
    api = os.environ.get('ERROR_API', "localhost:8080")
    params = os.environ.get("ERROR_API_PARAMS", "")
    conn = httplib.HTTPConnection(api)
    conn.request("GET", "/?lang=python&full=true"+ params)
    resp = conn.getresponse()
    print(resp.read())

sys.excepthook = webscale_errorhook
