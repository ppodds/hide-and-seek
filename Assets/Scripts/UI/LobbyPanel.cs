using UnityEngine;

namespace UI
{
    public class LobbyPanel : MonoBehaviour
    {
        [SerializeField] private Transform lobbyListTransform;
        [SerializeField] private GameObject lobbyItem;

        private void OnEnable()
        {
            var task = GameManager.Instance.GameTcpClient.GetLobbies();
            task.GetAwaiter().OnCompleted(() =>
            {
                foreach (var lobby in task.Result)
                {
                    var t = Instantiate(lobbyItem, lobbyListTransform);
                    t.GetComponent<LobbyItem>()
                        .SetText(lobby.Value.Lead, lobby.Value.CurrentPeople, lobby.Value.MaxPeople);
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
            Debug.Log("t");
            task.GetAwaiter().OnCompleted(() =>
            {
                Debug.Log("t2");
                if (task.Result == null)
                {
                    Debug.Log("Create lobby failed");
                    return;
                }

                Debug.Log("Create lobby success");
            });
        }
    }
}