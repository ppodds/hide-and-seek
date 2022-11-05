using System.Threading.Tasks;
using Protos;
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
            lead.SetText("Lead: " + Lobby.Lead.Id);
            people.SetText(Lobby.CurPeople + " / " + Lobby.MaxPeople);
        }

        private async Task JoinRoomTask()
        {
            var lobby = await GameManager.Instance.GameTcpClient.JoinLobby(Lobby);
            if (lobby == null)
            {
                Debug.Log("Join failed");
                return;
            }

            var result = await GameManager.Instance.GameUdpClient.ConnectLobby();
            if (!result.Success)
            {
                Debug.Log("Connect to lobby failed");
                return;
            }

            LobbyPanel.ShowPrepareRoom(lobby);
        }

        public void JoinRoom()
        {
            JoinRoomTask();
        }
    }
}