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

        public GameTcpClient(string host, int port)
        {
            _host = host;
            _port = port;
        }

        private async Task<byte[]> Rpc(byte procId, bool readResponse = true)
        {
            return await Rpc(procId, Array.Empty<byte>(), readResponse);
        }

        private async Task<byte[]> Rpc(byte procId, byte[] data, bool readResponse = true,
            CancellationToken token = default)
        {
            var client = new TcpClient();
            await client.ConnectAsync(_host, _port);
            await RpcCall(client, procId, data);
            var result = readResponse ? await ReadRpcResponse(client, token) : null;
            client.Close();
            client.Dispose();
            return result;
        }

        private static async Task RpcCall(TcpClient client, byte procId, byte[] data)
        {
            var stream = client.GetStream();
            var outputStream = new MemoryStream();
            await outputStream.WriteAsync(new[] { procId });
            await outputStream.WriteAsync(BitConverter.GetBytes(data.Length));
            await outputStream.WriteAsync(data);
            await stream.WriteAsync(outputStream.ToArray());
        }

        private static async Task<byte[]> ReadRpcResponse(TcpClient client, CancellationToken token = default)
        {
            var stream = client.GetStream();
            var buf = new byte[4];
            var n = await stream.ReadAsync(buf, token);
            token.ThrowIfCancellationRequested();
            if (n != buf.Length)
                throw new WrongProtocolException();
            var resLength = BitConverter.ToUInt32(buf);
            if (resLength == 0)
                return Array.Empty<byte>();
            buf = new byte[resLength];
            n = await stream.ReadAsync(buf, token);
            token.ThrowIfCancellationRequested();
            if (n != buf.Length)
                throw new WrongProtocolException();
            return buf;
        }

        public async Task<Player> Login()
        {
            return Player.Parser.ParseFrom(await Rpc(0));
        }

        public async Task<Lobbies> GetLobbies()
        {
            var buf = await Rpc(1);
            return buf == null ? null : Lobbies.Parser.ParseFrom(buf);
        }

        public async Task<Lobby> CreateLobby()
        {
            var player = new Player
            {
                Id = GameManager.Instance.PlayerID
            };
            var data = new CreateLobbyRequest
            {
                Lead = player
            };
            var outputStream = new MemoryStream();
            data.WriteTo(outputStream);
            var buf = await Rpc(2, outputStream.ToArray());
            return CreateLobbyResponse.Parser.ParseFrom(buf).Lobby;
        }

        public async Task<Lobby> JoinLobby(Lobby lobby)
        {
            var player = new Player
            {
                Id = GameManager.Instance.PlayerID
            };
            var data = new JoinLobbyRequest
            {
                Player = player,
                Lobby = lobby
            };
            var outputStream = new MemoryStream();
            data.WriteTo(outputStream);
            var buf = await Rpc(3, outputStream.ToArray());
            var res = JoinLobbyResponse.Parser.ParseFrom(buf);
            if (!res.Success) Debug.Log("Unable to join the lobby");
            return res.Lobby;
        }

        public async Task LeaveLobby(Lobby lobby)
        {
            var player = new Player
            {
                Id = GameManager.Instance.PlayerID
            };
            var data = new LeaveLobbyRequest
            {
                Player = player,
                Lobby = lobby
            };
            var outputStream = new MemoryStream();
            data.WriteTo(outputStream);
            var buf = await Rpc(4, outputStream.ToArray());
            var res = LeaveLobbyResponse.Parser.ParseFrom(buf);
            if (!res.Success) Debug.Log("Unable to leave the lobby");
        }

        public async Task Logout()
        {
            var player = new Player
            {
                Id = GameManager.Instance.PlayerID
            };
            var data = new LogoutRequest
            {
                Player = player
            };
            var outputStream = new MemoryStream();
            data.WriteTo(outputStream);
            await Rpc(5, outputStream.ToArray(), false);
        }

        public async Task<bool> StartGame(Lobby lobby)
        {
            var player = new Player
            {
                Id = GameManager.Instance.PlayerID
            };
            var data = new StartGameRequest
            {
                Player = player,
                Lobby = lobby
            };
            var outputStream = new MemoryStream();
            data.WriteTo(outputStream);
            var result = await Rpc(6, outputStream.ToArray());
            return StartGameResponse.Parser.ParseFrom(result).Success;
        }
    }
}