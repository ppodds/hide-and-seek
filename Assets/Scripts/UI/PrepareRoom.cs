using System;
using Protos;
using UnityEngine;

namespace UI
{
    public class PrepareRoom : MonoBehaviour
    {
        [SerializeField] private Transform playerListTransform;
        [SerializeField] private GameObject playerItem;
        [SerializeField] private GameObject lobbyList;

        private bool _trying;

        public Lobby Lobby { get; set; }

        private async void OnEnable()
        {
            if (Lobby != null)
                UpdatePrepareRoom();
            try
            {
                var inLobby = true;
                while (inLobby)
                {
                    var result = await GameManager.Instance.GameUdpClient.WaitLobbyBroadcast();
                    switch (result.Event)
                    {
                        case LobbyEvent.Join:
                        case LobbyEvent.Leave:
                            Lobby = result.Lobby;
                            UpdatePrepareRoom();
                            break;
                        case LobbyEvent.Destroy:
                            GameManager.Instance.DisconnectUdp();
                            lobbyList.SetActive(true);
                            gameObject.SetActive(false);
                            GameManager.Instance.toast.PushToast("Leave lobby success");
                            inLobby = false;
                            break;
                        case LobbyEvent.Start:
                            GameManager.Instance.StartGame(result.InitGame, Lobby);
                            inLobby = false;
                            break;
                        default:
                            throw new ArgumentOutOfRangeException();
                    }
                }
            }
            catch (ObjectDisposedException e) // player has already leave the lobby
            {
            }
        }

        private void OnDisable()
        {
            ClearPlayerList();
        }

        private void ClearPlayerList()
        {
            for (var i = 0; i < playerListTransform.childCount; i++)
                Destroy(playerListTransform.GetChild(i).gameObject);
        }

        public void UpdatePrepareRoom()
        {
            ClearPlayerList();
            foreach (var player in Lobby.Players)
            {
                var t = Instantiate(playerItem, playerListTransform);
                t.GetComponent<PlayerItem>().SetText(Lobby, player);
            }
        }

        public void BackToLobbyList()
        {
            if (_trying)
                return;
            _trying = true;
            GameManager.Instance.DisconnectUdp();
            var task = GameManager.Instance.GameTcpClient.LeaveLobby(Lobby);
            task.GetAwaiter().OnCompleted(() =>
            {
                _trying = false;
                lobbyList.SetActive(true);
                gameObject.SetActive(false);
                GameManager.Instance.toast.PushToast("Leave lobby success");
            });
        }

        public void StartGame()
        {
            var task = GameManager.Instance.GameTcpClient.StartGame(Lobby);
            task.GetAwaiter().OnCompleted(() =>
            {
                if (!task.Result)
                    GameManager.Instance.toast.PushToast("Start game failed");
            });
        }
    }
}