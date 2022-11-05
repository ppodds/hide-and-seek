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

        private async Task RpcCall(byte procId, byte[] data)
        {
            var outputStream = new MemoryStream();
            await outputStream.WriteAsync(new[] { procId });
            await outputStream.WriteAsync(BitConverter.GetBytes(data.Length));
            var t = outputStream.ToArray();
            await _udpClient.SendAsync(t, t.Length);
            await _udpClient.SendAsync(data, data.Length);
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
                Id = GameManager.Instance.ID
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
    }
}