using System;
using System.Collections.Generic;
using System.Linq;
using System.Net.Sockets;
using System.Text;
using System.Threading.Tasks;
using Newtonsoft.Json;
using Server;
using UnityEngine;

namespace IO.Net
{
    public class GameTcpClient
    {
        private readonly string _host;
        private readonly int _port;
        private readonly TcpClient _tcpClient;

        public GameTcpClient(string host, int port)
        {
            _tcpClient = new TcpClient();
            _host = host;
            _port = port;
        }

        public async Task<uint> Login()
        {
            await _tcpClient.ConnectAsync(_host, _port);
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
            buf = new byte[1024];
            var n = await stream.ReadAsync(buf);
            var lobby = JsonConvert.DeserializeObject<Lobby>(Encoding.UTF8.GetString(buf));
            return lobby;
        }

        public async Task<Lobby> JoinLobby(uint id)
        {
            var stream = _tcpClient.GetStream();
            var buf = new byte[] { 3 }.Concat(BitConverter.GetBytes(GameManager.Instance.ID))
                .Concat(BitConverter.GetBytes(id)).ToArray();
            await stream.WriteAsync(buf);
            buf = new byte[1024];
            var n = await stream.ReadAsync(buf);
            if (n == 1 && buf[0] == 0)
            {
                Debug.Log("Unable to join the lobby");
                return null;
            }

            var lobby = JsonConvert.DeserializeObject<Lobby>(Encoding.UTF8.GetString(buf));
            return lobby;
        }

        public async Task<Lobby> LeaveLobby(uint id)
        {
            var stream = _tcpClient.GetStream();
            var buf = new byte[] { 4 }.Concat(BitConverter.GetBytes(GameManager.Instance.ID))
                .Concat(BitConverter.GetBytes(id)).ToArray();
            await stream.WriteAsync(buf);
            buf = new byte[1024];
            var n = await stream.ReadAsync(buf);
            if (n == 1 && buf[0] == 0)
            {
                Debug.Log("Unable to leave the lobby");
                return null;
            }

            var lobby = JsonConvert.DeserializeObject<Lobby>(Encoding.UTF8.GetString(buf));
            return lobby;
        }
    }
}