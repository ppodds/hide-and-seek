using System;
using System.Net.Sockets;
using System.Threading.Tasks;
using UnityEngine;

namespace IO.Net
{
    public class GameTcpClient
    {
        private readonly TcpClient _tcpClient;

        public GameTcpClient(string host, int port)
        {
            _tcpClient = new TcpClient();
            _tcpClient.Connect(host, port);
        }

        public async Task<uint> Login()
        {
            var stream = _tcpClient.GetStream();
            var buf = new byte[] { 0 };
            await stream.WriteAsync(buf);
            buf = new byte[4];
            var n = await stream.ReadAsync(buf);
            if (n != 4)
                throw new WrongProtocolException();
            return BitConverter.ToUInt32(buf);
        }
    }
}