using Server;
using UnityEngine;

namespace UI
{
    public class LobbyPanel : MonoBehaviour
    {
        [SerializeField] private Transform lobbyListTransform;
        [SerializeField] private GameObject lobbyItem;
        [SerializeField] private GameObject prepareRoom;

        private void OnEnable()
        {
            var task = GameManager.Instance.GameTcpClient.GetLobbies();
            task.GetAwaiter().OnCompleted(() =>
            {
                foreach (var lobby in task.Result)
                {
                    var t = Instantiate(lobbyItem, lobbyListTransform).GetComponent<LobbyItem>();
                    t.Lobby = lobby.Value;
                    t.LobbyPanel = this;
                    t.UpdateText();
                }

                Debug.Log("Load lobbies success");
            });
        }

        private void OnDisable()
        {
            for (var i = 0; i < lobbyListTransform.childCount; i++) Destroy(lobbyListTransform.GetChild(i).gameObject);
        }

        public void CreateLobby()
        {
            var task = GameManager.Instance.GameTcpClient.CreateLobby();
            task.GetAwaiter().OnCompleted(() =>
            {
                if (task.Result == null)
                {
                    Debug.Log("Create lobby failed");
                    return;
                }

                ShowPrepareRoom(task.Result);
                Debug.Log("Create lobby success");
            });
        }

        public void ShowPrepareRoom(Lobby lobby)
        {
            prepareRoom.SetActive(true);
            var t = prepareRoom.GetComponent<PrepareRoom>();
            t.Lobby = lobby;
            t.UpdatePrepareRoom();
            gameObject.SetActive(false);
        }
    }
}