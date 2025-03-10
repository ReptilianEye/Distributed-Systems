import socket;

serverIP = "127.0.0.1"
serverPort = 9011
msg = "Żółta gęś!"

print('PYTHON UDP CLIENT')
client = socket.socket(socket.AF_INET, socket.SOCK_DGRAM)
client.sendto(bytes(msg, 'utf8'), (serverIP, serverPort))




