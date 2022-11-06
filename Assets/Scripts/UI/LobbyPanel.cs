using System.Threading.Tasks;
using Protos;
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
                if (task.Result == null)
                    return;
                foreach (var lobby in task.Result.Lobbies_)
                {
                    var t = Instantiate(lobbyItem, lobbyListTransform).GetComponent<LobbyItem>();
                    t.Lobby = lobby.Value;
                    t.UpdateText();
                }

                Debug.Log("Load lobbies success");
            });
        }

        private void OnDisable()
        {
            for (var i = 0; i < lobbyListTransform.childCount; i++) Destroy(lobbyListTransform.GetChild(i).gameObject);
        }

        private async Task Join()
        {
            var lobby = await GameManager.Instance.GameTcpClient.CreateLobby();
            if (lobby == null)
            {
                Debug.Log("Create lobby failed");
                return;
            }

            await GameManager.Instance.ConnectToLobby(lobby);
        }

        public void CreateLobby()
        {
            Join();
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