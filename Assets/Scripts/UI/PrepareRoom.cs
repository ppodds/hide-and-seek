using System;
using System.Threading;
using Protos;
using UnityEngine;

namespace UI
{
    public class PrepareRoom : MonoBehaviour
    {
        [SerializeField] private Transform playerListTransform;
        [SerializeField] private GameObject playerItem;
        [SerializeField] private GameObject lobbyList;

        private CancellationTokenSource _tokenSource;

        private bool _trying;

        public Lobby Lobby { get; set; }

        private async void OnEnable()
        {
            _tokenSource = new CancellationTokenSource();
            try
            {
                while (true)
                {
                    var result = await GameManager.Instance.GameUdpClient.WaitLobbyBroadcast(_tokenSource.Token);
                    switch (result.Event)
                    {
                        case LobbyEvent.Join:
                        case LobbyEvent.Leave:
                            Lobby = result.Lobby;
                            UpdatePrepareRoom();
                            break;
                        case LobbyEvent.Destroy:
                            BackToLobbyList();
                            break;
                        case LobbyEvent.Start:
                            break;
                        default:
                            throw new ArgumentOutOfRangeException();
                    }
                }
            }
            catch (OperationCanceledException e)
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
            _tokenSource.Cancel();
            _tokenSource.Dispose();
            var task = GameManager.Instance.GameTcpClient.LeaveLobby(Lobby);
            task.GetAwaiter().OnCompleted(() =>
            {
                _trying = false;
                lobbyList.SetActive(true);
                gameObject.SetActive(false);
                Debug.Log("Leave lobby success");
            });
        }

        public void StartGame()
        {
        }
    }
}