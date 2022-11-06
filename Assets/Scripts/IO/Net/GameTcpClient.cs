using System;
using System.IO;
using System.Net.Sockets;
using System.Threading;
using System.Threading.Tasks;
using Google.Protobuf;
using Protos;
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

        private async Task RpcCall(byte procId)
        {
            await RpcCall(procId, Array.Empty<byte>());
        }

        private async Task RpcCall(byte procId, byte[] data)
        {
            var stream = _tcpClient.GetStream();
            var outputStream = new MemoryStream();
            await outputStream.WriteAsync(new[] { procId });
            await outputStream.WriteAsync(BitConverter.GetBytes(data.Length));
            await outputStream.WriteAsync(data);
            await stream.WriteAsync(outputStream.ToArray());
        }

        private async Task<byte[]> ReadRpcResponse(CancellationToken token = default)
        {
            var stream = _tcpClient.GetStream();
            var buf = new byte[4];
            var n = await stream.ReadAsync(buf, token);
            token.ThrowIfCancellationRequested();
            if (n != buf.Length)
                throw new WrongProtocolException();
            var resLength = BitConverter.ToUInt32(buf);
            if (resLength == 0)
                return null;
            buf = new byte[resLength];
            n = await stream.ReadAsync(buf, token);
            token.ThrowIfCancellationRequested();
            if (n != buf.Length)
                throw new WrongProtocolException();
            return buf;
        }

        public async Task<Player> Login()
        {
            await _tcpClient.ConnectAsync(_host, _port);
            await RpcCall(0);
            var data = await ReadRpcResponse();
            return Player.Parser.ParseFrom(data);
        }

        public async Task<Lobbies> GetLobbies()
        {
            await RpcCall(1);
            var buf = await ReadRpcResponse();
            return buf == null ? null : Lobbies.Parser.ParseFrom(buf);
        }

        public async Task<Lobby> CreateLobby()
        {
            var player = new Player
            {
                Id = GameManager.Instance.ID
            };
            var data = new CreateLobbyRequest
            {
                Lead = player
            };
            var outputStream = new MemoryStream();
            data.WriteTo(outputStream);
            await RpcCall(2, outputStream.ToArray());
            var buf = await ReadRpcResponse();
            return CreateLobbyResponse.Parser.ParseFrom(buf).Lobby;
        }

        public async Task<Lobby> JoinLobby(Lobby lobby)
        {
            var player = new Player
            {
                Id = GameManager.Instance.ID
            };
            var data = new JoinLobbyRequest
            {
                Player = player,
                Lobby = lobby
            };
            var outputStream = new MemoryStream();
            data.WriteTo(outputStream);
            await RpcCall(3, outputStream.ToArray());
            var buf = await ReadRpcResponse();
            var res = JoinLobbyResponse.Parser.ParseFrom(buf);
            if (!res.Success) Debug.Log("Unable to join the lobby");
            return res.Lobby;
        }

        public async Task LeaveLobby(Lobby lobby)
        {
            var player = new Player
            {
                Id = GameManager.Instance.ID
            };
            var data = new LeaveLobbyRequest
            {
                Player = player,
                Lobby = lobby
            };
            var outputStream = new MemoryStream();
            data.WriteTo(outputStream);
            await RpcCall(4, outputStream.ToArray());
            var buf = await ReadRpcResponse();
            var res = LeaveLobbyResponse.Parser.ParseFrom(buf);
            if (!res.Success) Debug.Log("Unable to leave the lobby");
        }

        public async Task Logout()
        {
            var player = new Player
            {
                Id = GameManager.Instance.ID
            };
            var data = new LogoutRequest
            {
                Player = player
            };
            var outputStream = new MemoryStream();
            data.WriteTo(outputStream);
            await RpcCall(5, outputStream.ToArray());
        }
    }
}