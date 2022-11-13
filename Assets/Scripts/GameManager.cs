using System;
using System.Net.Sockets;
using System.Threading;
using System.Threading.Tasks;
using IO.Net;
using Protos;
using SUPERCharacter;
using UI;
using Unity.VisualScripting;
using UnityEngine;
using UnityEngine.SceneManagement;
using Vector3 = UnityEngine.Vector3;

public struct Server
{
    public string Host;
    public int TcpPort;
    public int UdpPort;
}

public class GameManager : MonoBehaviour
{
    [SerializeField] private LobbyPanel lobbyPanel;
    [SerializeField] private GameObject menuUI;

    public Transform PlayersParent;
    public uint PlayerID { get; private set; }

    public GameState GameState { get; private set; }

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
            PlayerID = (await GameTcpClient.Login()).Id;
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

    public void Logout()
    {
        var task = GameTcpClient.Logout();
        task.GetAwaiter().OnCompleted(Application.Quit);
    }

    private async Task<bool> ConnectToGame()
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

        var result = await GameUdpClient.ConnectGame();
        if (!result.Success)
        {
            Debug.Log("Connect to game failed");
            return false;
        }

        SceneManager.LoadScene("Demo");
        Debug.Log("Join game success");
        return true;
    }

    public async Task StartGame(InitGame initGame, Lobby lobby)
    {
        DisconnectUdp();
        await ConnectToGame();
        menuUI.SetActive(false);
        GameState = new GameState(initGame.Game.Id);
        foreach (var pair in initGame.Players)
        {
            var o = pair.Value.Character.Type == CharacterType.Player
                ? Instantiate(Resources.Load("Prefabs/Player/Human"), PlayersParent)
                : Instantiate(Resources.Load("Prefabs/Player/Ghost"), PlayersParent);
            var netO = o.GetComponent<NetObject>();
            netO.IsRemote = pair.Key != PlayerID;
            netO.PlayerId = pair.Key;
            netO.transform.position = new Vector3(pair.Value.Character.Pos.X, pair.Value.Character.Pos.Y,
                pair.Value.Character.Pos.Z);
            if (pair.Key == PlayerID)
            {
                var charController = netO.GetComponent<SUPERCharacterAIO>();
                charController.enabled = true;
                charController.playerCamera.gameObject.SetActive(true);
            }

            GameState.Players.Add(pair.Key, new PlayerState
            {
                Player = pair.Value,
                PlayerObject = netO
            });
        }

        var handlePlayerMoveThread = new Thread(() => { HandlePlayerMove(); });
        handlePlayerMoveThread.Start();
    }

    public async Task HandlePlayerMove()
    {
        try
        {
            while (true)
            {
                var result = await GameUdpClient.WaitPlayerUpdateBroadcast();
                var player = result.Player;
                if (player.Player.Id == PlayerID)
                    continue;
                GameState.Players[player.Player.Id].Player = result.Player;
            }
        }
        catch (ObjectDisposedException e) // player has already leave the lobby
        {
        }
    }
}