using System;
using System.IO;
using System.Net.Sockets;
using System.Threading;
using System.Threading.Tasks;
using Google.Protobuf;
using Protos;

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

        public void Connect()
        {
            _udpClient.Connect(_host, _port);
        }

        public void Disconnect()
        {
            _udpClient.Close();
        }

        private async Task RpcCall(byte procId, byte[] data)
        {
            var outputStream = new MemoryStream();
            await outputStream.WriteAsync(new[] { procId });
            await outputStream.WriteAsync(BitConverter.GetBytes(data.Length));
            await outputStream.WriteAsync(data);
            var t = outputStream.ToArray();
            await _udpClient.SendAsync(t, t.Length);
        }

        private async Task<byte[]> ReadRpcResponse()
        {
            return (await _udpClient.ReceiveAsync()).Buffer;
        }

        private async Task<byte[]> ReadBroadcast(CancellationToken token = default)
        {
            var result = await _udpClient.ReceiveAsync();
            token.ThrowIfCancellationRequested();
            return result.Buffer;
        }

        public async Task<ConnectLobbyResponse> ConnectLobby()
        {
            var player = new Player
            {
                Id = GameManager.Instance.PlayerID
            };
            var data = new ConnectLobbyRequest
            {
                Player = player
            };
            var outputStream = new MemoryStream();
            data.WriteTo(outputStream);
            await RpcCall(0, outputStream.ToArray());
            var buf = await ReadRpcResponse();
            return ConnectLobbyResponse.Parser.ParseFrom(buf);
        }

        public async Task<LobbyBroadcast> WaitLobbyBroadcast(CancellationToken token = default)
        {
            var res = await ReadBroadcast(token);
            return LobbyBroadcast.Parser.ParseFrom(res);
        }

        public async Task<ConnectGameResponse> ConnectGame()
        {
            var player = new Player
            {
                Id = GameManager.Instance.PlayerID
            };
            var data = new ConnectGameRequest
            {
                Player = player
            };
            var outputStream = new MemoryStream();
            data.WriteTo(outputStream);
            await RpcCall(1, outputStream.ToArray());
            var buf = await ReadRpcResponse();
            return ConnectGameResponse.Parser.ParseFrom(buf);
        }

        public async Task<GameBroadcast> WaitGameBroadcast(CancellationToken token = default)
        {
            var res = await ReadBroadcast(token);
            return GameBroadcast.Parser.ParseFrom(res);
        }

        public async Task UpdatePlayer(Character character)
        {
            var data = new UpdatePlayerRequest
            {
                Game = new Game
                {
                    Id = GameManager.Instance.GameState.Id
                },
                Player = new GamePlayer
                {
                    Player = new Player
                    {
                        Id = GameManager.Instance.PlayerID
                    },
                    Character = character
                }
            };
            var outputStream = new MemoryStream();
            data.WriteTo(outputStream);
            await RpcCall(2, outputStream.ToArray());
        }

        public async Task<UpdatePlayerBroadcast> WaitPlayerUpdateBroadcast(CancellationToken token = default)
        {
            var res = await ReadBroadcast(token);
            return UpdatePlayerBroadcast.Parser.ParseFrom(res);
        }
    }
}