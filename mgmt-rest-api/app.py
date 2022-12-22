from flask import Flask, request
from redis import Redis

app = Flask(__name__)

@app.route('/api/v1/health', methods=['GET'])
def health():
    return 'OK', 200

@app.route('/dns', methods=['GET'])
def dnsGET():
    """
    Makes a call to the redis server to get the dns info
    """

    # get query params

    domainName = request.args.get('domainName')

    # Get the dns info from redis
    r = Redis(host='redis', port=6379, db=0)
    ip = r.get(domainName)

    # Return the dns info
    return ip, 200

@app.route('/dns', methods=['POST'])
def dnsPOST():
    """
    Makes a call to the redis server to set the dns info
    """

    # ip and domainName from request body
    ip = request.json['ip']
    domainName = request.json['domainName']

    # Set the dns info in redis
    r = Redis(host='redis', port=6379, db=0)
    r.set(domainName, ip)

    # Return the dns info
    return ip, 200

if __name__ == '__main__':
    app.run(host="0.0.0.0", port=5000)
