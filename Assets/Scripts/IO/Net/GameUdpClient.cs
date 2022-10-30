using System.Net.Sockets;

namespace IO.Net
{
    public class GameUdpClient
    {
        private readonly UdpClient _udpClient;

        public GameUdpClient(string host, int port)
        {
            _udpClient = new UdpClient();
            _udpClient.Connect(host, port);
        }
    }
}