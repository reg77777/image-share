from flask import Flask, request
from dotenv import load_dotenv
import hmac
import os
import hashlib
from multiprocessing import Process
import subprocess

load_dotenv('.env')
webhook_secret = os.environ['WEBHOOK_SECRET']

app = Flask(__name__)

def verify(payload, secret, signature):
    if not signature:
        return False
    hash_object = hmac.new(bytes(secret.encode(
        'utf-8')), msg=payload, digestmod=hashlib.sha256)
    expected_signature = "sha256="+hash_object.hexdigest()
    if not hmac.compare_digest(expected_signature, signature):
        return False
    return True

def deploy_process(branch):
    cmd = ['git','pull','--all']
    proc = subprocess.run(cmd)

    cmd = ['git','checkout',branch]
    proc = subprocess.run(cmd)

    cmd = ['docker-compose','down']
    proc = subprocess.run(cmd)

    cmd = ['docker-compose','up','-d']
    proc = subprocess.run(cmd)

@app.route('/',methods=['POST'])
def deploy():
    print('aaaa')
    if verify(request.get_data(), webhook_secret, request.headers.get('X-Hub-Signature-256')):
        ref = request.json['ref']
        branch = ref.split('/')[-1]
        print('branch: ',branch)
        p=Process(target=deploy_process,args=(branch,))
        p.start()
        return ''
    else:
        print('invalid signature')
        return

if __name__ == '__main__':
    app.run(debug=False, host='0.0.0.0', port=3000)
