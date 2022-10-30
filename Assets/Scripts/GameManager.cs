﻿using System.Net.Sockets;
using System.Threading.Tasks;
using IO.Net;
using UnityEngine;

public class GameManager : MonoBehaviour
{
    public uint ID { get; private set; }

    public GameTcpClient GameTcpClient { get; private set; }

    public GameUdpClient GameUdpClient { get; private set; }

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
            GameTcpClient = new GameTcpClient(host, tcpPort);
            GameUdpClient = new GameUdpClient(host, udpPort);
        }
        catch (SocketException e)
        {
            return false;
        }

        ID = await GameTcpClient.Login();
        return true;
    }
}