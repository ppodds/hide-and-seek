using Server;
using TMPro;
using UnityEngine;

namespace UI
{
    public class LobbyItem : MonoBehaviour
    {
        [SerializeField] private TMP_Text lead;
        [SerializeField] private TMP_Text people;

        public LobbyPanel LobbyPanel { get; set; }
        public Lobby Lobby { get; set; }

        public void UpdateText()
        {
            lead.SetText("Lead: " + Lobby.Lead.ID);
            people.SetText(Lobby.CurrentPeople + " / " + Lobby.MaxPeople);
        }

        public void JoinRoom()
        {
            var task = GameManager.Instance.GameTcpClient.JoinLobby(Lobby.ID);
            task.GetAwaiter().OnCompleted(() =>
            {
                if (task.Result == null)
                {
                    Debug.Log("Join failed");
                    return;
                }

                LobbyPanel.ShowPrepareRoom(task.Result);
            });
        }
    }
}