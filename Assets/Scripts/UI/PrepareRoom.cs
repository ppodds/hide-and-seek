using Server;
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

        private void OnDisable()
        {
            for (var i = 0; i < playerListTransform.childCount; i++)
                Destroy(playerListTransform.GetChild(i).gameObject);
        }

        public void UpdatePrepareRoom()
        {
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
            var task = GameManager.Instance.GameTcpClient.LeaveLobby(Lobby.ID);
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