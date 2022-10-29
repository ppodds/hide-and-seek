using System;
using System.Net.Sockets;
using System.Threading.Tasks;
using IO;
using IO.Net;
using UnityEngine;

public class GameManager : MonoBehaviour
{
    private GameTcpClient _gameTcpClient;
    private GameUdpClient _gameUdpClient;

    private uint _id;

    public static GameManager Instance { get; private set; }

    private void Awake()
    {
        if (Instance != null)
        {
            Destroy(gameObject);
            return;
        }

        Instance = this;
        DontDestroyOnLoad(gameObject);
    }
    
    public async Task<bool> ConnectToServer(string host, int tcpPort, int udpPort)
    {
        try
        {
            _gameTcpClient = new GameTcpClient(host, tcpPort);
            _gameUdpClient = new GameUdpClient(host, udpPort);
        }
        catch (SocketException e)
        {
            return false;
        }

        _id = await _gameTcpClient.Login();
        return true;
    }
}