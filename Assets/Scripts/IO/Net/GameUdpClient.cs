using System.Net.Sockets;

namespace IO.Net
{
    public class GameUdpClient
    {
        private readonly string _host;
        private readonly int _port;
        private readonly UdpClient _udpClient;

        public GameUdpClient(string host, int port)
        {
            _udpClient = new UdpClient();
            _host = host;
            _port = port;
        }

        public void Login()
        {
            _udpClient.Connect(_host, _port);
        }
    }
}