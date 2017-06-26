#!/usr/bin/env python

from lib.google.cloud import pubsub
import subprocess

def message_loop(handler):

    client = pubsub.Client.from_service_account_json('account.json',  project='kodicloud-169614')
    topic = client.topic('AirCon')
    subscription = topic.subscription('Test')

    while True:
        for ack_id, message in subscription.pull():
            try:
                handler(message.data)
            except Exception as e:
                print e
            finally:
                subscription.acknowledge([ack_id])

def send_AC_command(command):

    if command != 'off':
        try:
            if int(command) not in xrange(16, 31):
                raise ValueError()
        except ValueError:
            raise ValueError('invalid A/C command: ' + command)

    for _ in xrange(1):
        print 'executing command: ' + command
        subprocess.check_call(['irsend', 'SEND_ONCE', 'ac', command])


if __name__ == '__main__':
    message_loop(send_AC_command)
