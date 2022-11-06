﻿using System.Net.Sockets;
using System.Threading.Tasks;
using IO.Net;
using Protos;
using UI;
using UnityEngine;

public struct Server
{
    public string Host;
    public int TcpPort;
    public int UdpPort;
}

public class GameManager : MonoBehaviour
{
    [SerializeField] private LobbyPanel lobbyPanel;
    public uint ID { get; private set; }

    public GameTcpClient GameTcpClient { get; private set; }

    public GameUdpClient GameUdpClient { get; private set; }

    public static GameManager Instance { get; private set; }

    public Server Server { get; set; }

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

    public async Task<bool> ConnectToServer()
    {
        GameTcpClient = new GameTcpClient(Server.Host, Server.TcpPort);
        try
        {
            ID = (await GameTcpClient.Login()).Id;
        }
        catch (SocketException e)
        {
            return false;
        }

        return true;
    }

    public async Task<bool> ConnectToLobby(Lobby lobby)
    {
        GameUdpClient = new GameUdpClient(Server.Host, Server.UdpPort);
        try
        {
            GameUdpClient.Connect();
        }
        catch (SocketException e)
        {
            return false;
        }

        var result = await GameUdpClient.ConnectLobby();
        if (!result.Success)
        {
            Debug.Log("Connect to lobby failed");
            return false;
        }

        lobbyPanel.ShowPrepareRoom(lobby);
        Debug.Log("Join lobby success");
        return true;
    }

    public void DisconnectUdp()
    {
        GameUdpClient.Disconnect();
        GameUdpClient = null;
    }
}