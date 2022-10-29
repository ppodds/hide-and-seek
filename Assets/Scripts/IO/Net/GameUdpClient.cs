using System;
using System.Net.Sockets;
using System.Threading.Tasks;
using UnityEngine;

namespace IO.Net
{
    public class GameUdpClient
    {
        private UdpClient _udpClient;

        public GameUdpClient(string host, int port)
        {
            _udpClient = new UdpClient();
            _udpClient.Connect(host, port);
        }
    }
}