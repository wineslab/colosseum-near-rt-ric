import socket


if __name__ == '__main__':

    client = socket.socket(socket.AF_INET, socket.SOCK_STREAM)
    client.connect(('127.0.0.1', 7000))

    msg = 'terminate'
    bytes_num = client.send(msg.encode('utf-8'))
    print('Socket sent ' + str(bytes_num) + ' bytes')

    client.close()
