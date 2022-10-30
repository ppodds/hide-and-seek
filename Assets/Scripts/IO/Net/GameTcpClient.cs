using System;
using System.Collections.Generic;
using System.Linq;
using System.Net.Sockets;
using System.Text;
using System.Threading.Tasks;
using Newtonsoft.Json;
using Server;

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

        public async Task<Dictionary<uint, Lobby>> GetLobbies()
        {
            var stream = _tcpClient.GetStream();
            var buf = new byte[] { 1 };
            await stream.WriteAsync(buf);
            buf = new byte[4096];
            var n = await stream.ReadAsync(buf);
            var lobbies = JsonConvert.DeserializeObject<Dictionary<uint, Lobby>>(Encoding.UTF8.GetString(buf));
            return lobbies;
        }

        public async Task<Lobby> CreateLobby()
        {
            var stream = _tcpClient.GetStream();
            var buf = new byte[] { 2 }.Concat(BitConverter.GetBytes(GameManager.Instance.ID)).ToArray();
            await stream.WriteAsync(buf);
            buf = new byte[4096];
            var n = await stream.ReadAsync(buf);
            var lobby = JsonConvert.DeserializeObject<Lobby>(Encoding.UTF8.GetString(buf));
            return lobby;
        }
    }
}