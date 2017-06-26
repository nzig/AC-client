import client
import xbmc

if __name__ == '__main__':
    monitor = xbmc.Monitor()
    def abort():
        return monitor.abortRequested()

    client.message_loop(client.send_AC_command, abort)
